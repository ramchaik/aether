import { fetchApiUrl } from "./awsApiGateway";

export const getProjectDomain = async (projectId: string) => {
  const baseDomain = await fetchApiUrl();
  return `${baseDomain}/${projectId}/index.html`;
};
