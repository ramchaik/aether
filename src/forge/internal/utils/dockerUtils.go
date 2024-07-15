package utils

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

// buildImage builds a Docker image using the provided Dockerfile and options.
func buildImage(ctx context.Context, cli *client.Client, dockerfilePath, repoURL, buildCommand, imageName string) (io.ReadCloser, error) {
	dockerfileContent, err := os.ReadFile(dockerfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Dockerfile: %w", err)
	}

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	hdr := &tar.Header{
		Name: "Dockerfile",
		Mode: 0600,
		Size: int64(len(dockerfileContent)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return nil, fmt.Errorf("failed to write tar header: %w", err)
	}
	if _, err := tw.Write(dockerfileContent); err != nil {
		return nil, fmt.Errorf("failed to write Dockerfile to tar archive: %w", err)
	}

	buildContext := ctx
	imageBuildResponse, err := cli.ImageBuild(buildContext, buf, types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		BuildArgs:  map[string]*string{"REPO_URL": &repoURL, "BUILD_COMMAND": &buildCommand},
		Tags:       []string{imageName},
		Remove:     true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build the image: %w", err)
	}

	// Read the response from the image build and log it
	log.Println("Reading image build response...")
	if _, err := io.Copy(os.Stdout, imageBuildResponse.Body); err != nil {
		log.Printf("Failed to read image build response: %v\n", err)
		return nil, err // Consider whether you want to return here or continue execution
	}

	return imageBuildResponse.Body, nil
}

// copyBuildOutput copies the build output from the container to the host.
func copyBuildOutput(ctx context.Context, cli *client.Client, imageName, currentDir string) error {
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   []string{"sh", "-c", "cp -r /build/* /app/build-output/"},
	}, &container.HostConfig{
		Binds: []string{
			currentDir + ":/app",
		},
	}, nil, nil, "")
	if err != nil {
		return fmt.Errorf("failed to create the container: %w", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start the container: %w", err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("error waiting for container: %w", err)
		}
	case <-statusCh:
	}

	return nil
}

func removeBuildDirectory(path string) error {
	return os.RemoveAll(path)
}

func pruneDockerImages(ctx context.Context, cli *client.Client) error {
	// pruning dangling images
	_, err := cli.ImagesPrune(ctx, filters.NewArgs())
	if err != nil {
		return fmt.Errorf("failed to prune Docker images: %w", err)
	}

	log.Println("Docker images pruned successfully")
	return nil
}

// BuildProject builds a project and returns the Docker client, build directory, and image name.
func BuildProject(ctx context.Context, repoURL, buildCommand string) (*client.Client, string, string, error) {
	uuid := uuid.New().String()
	imageName := "aether-build-" + uuid

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to get current directory: %w", err)
	}

	buildDir := filepath.Join(currentDir, "build-output")
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return nil, "", "", fmt.Errorf("failed to create build directory: %w", err)
	}

	dockerfilePath := filepath.Join(currentDir, "internal", "utils", "secure-build.dockerfile")

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create Docker client: %w", err)
	}

	_, err = buildImage(ctx, cli, dockerfilePath, repoURL, buildCommand, imageName)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed during image build: %w", err)
	}

	if err := copyBuildOutput(ctx, cli, imageName, currentDir); err != nil {
		return nil, "", "", fmt.Errorf("failed to copy build output: %w", err)
	}

	return cli, buildDir, imageName, nil
}

// Cleanup performs cleanup actions after a build project.
func Cleanup(ctx context.Context, cli *client.Client, buildDir string, imageName string) error {
	// Remove the build directory
	if err := removeBuildDirectory(buildDir); err != nil {
		return fmt.Errorf("failed to remove build directory: %w", err)
	}

	// Prune Docker images
	if err := pruneDockerImages(ctx, cli); err != nil {
		return fmt.Errorf("failed to prune Docker images: %w", err)
	}

	fmt.Println("Cleanup completed successfully")
	return nil
}
