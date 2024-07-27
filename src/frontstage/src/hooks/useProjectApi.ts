import { useMutation, useQuery, useQueryClient } from "react-query";
import { useAuth } from "@clerk/nextjs";

const ALL_PROJECTS_QUERY_KEY = "allProjects";
const PROJECT_QUERY_KEY = "project";
const PROJECT_DEPLOY_QUERY_KEY = "projectDeploy";
const PROJECT_LOGS_QUERY_KEY = "projectLogs";

interface ProjectData {
  name: string;
  repositoryUrl: string;
  customDomain?: string;
  buildCommand?: string;
}

const createProject = async (
  projectData: ProjectData,
  token: string | null
) => {
  const response = await fetch(`/api/project`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(projectData),
  });

  if (!response.ok) {
    throw new Error("Something went wrong");
  }

  return response.json();
};

const fetchAllProjects = async (token: string | null) => {
  const response = await fetch(`/api/project`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    throw new Error("Something went wrong");
  }

  return response.json();
};

const fetchProject = async (projectId: string, token: string | null) => {
  const response = await fetch(`/api/project/${projectId}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    throw new Error("Something went wrong");
  }

  return response.json();
};

export const deployProject = async (
  projectId: string,
  token: string | null
) => {
  const response = await fetch(`/api/project/${projectId}/deploy`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({}),
  });

  if (!response.ok) {
    throw new Error("Something went wrong");
  }

  return response.json();
};

export const fetchProjectLogs = async (
  projectId: string,
  token: string | null
) => {
  const response = await fetch(`/api/project/${projectId}/log`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    throw new Error("Something went wrong");
  }

  return response.json();
};

export function useCreateProject() {
  const { getToken } = useAuth();
  const queryClient = useQueryClient();

  return useMutation(
    async (projectData: ProjectData) => {
      const token = await getToken();
      return createProject(projectData, token);
    },
    {
      onSuccess: () => {
        // Invalidate and refetch the all projects query
        queryClient.invalidateQueries(ALL_PROJECTS_QUERY_KEY);
      },
    }
  );
}

export function useFetchAllProjects<T>() {
  const { getToken } = useAuth();

  return useQuery<T>({
    queryKey: ALL_PROJECTS_QUERY_KEY,
    queryFn: async () => {
      const token = await getToken();
      return fetchAllProjects(token);
    },
  });
}

export function useFetchProject<T>(projectId: string, options = {}) {
  const { getToken } = useAuth();

  return useQuery<T>({
    queryKey: [PROJECT_QUERY_KEY, projectId],
    queryFn: async () => {
      const token = await getToken();
      return fetchProject(projectId, token);
    },
    // Only fetch when projectId is available
    enabled: !!projectId,
    ...options,
  });
}

export function useDeployProject() {
  const { getToken } = useAuth();
  const queryClient = useQueryClient();

  return useMutation(
    async (projectId: string) => {
      const token = await getToken();
      return deployProject(projectId, token);
    },
    {
      onSuccess: (_, projectId) => {
        // Invalidate and refetch the project query
        queryClient.invalidateQueries([PROJECT_QUERY_KEY, projectId]);
      },
    }
  );
}

export function useFetchProjectLogs<T>(projectId: string, options = {}) {
  const { getToken } = useAuth();

  return useQuery<T>({
    queryKey: [PROJECT_LOGS_QUERY_KEY, projectId],
    queryFn: async () => {
      const token = await getToken();
      return fetchProjectLogs(projectId, token);
    },
    // Only fetch when projectId is available
    enabled: !!projectId,
    ...options,
  });
}
