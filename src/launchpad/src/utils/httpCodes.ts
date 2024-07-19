export const HTTP_CODES = {
  OK: 200,
  CREATED: 201,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  CONFLICT: 409,
  INTERNAL_SERVER_ERROR: 500,
} as const;

export const ERROR_MESSAGES = {
  INVALID_INPUT: "Invalid input",
  DUPLICATE_CUSTOM_DOMAIN:
    "Custom domain already exits. Please use another domain",
  PROJECT_CREATE_FAILED: "Failed to create project",
  PROJECT_NOT_FOUND: "Project not found",
  DEPLOYMENT_FAILED: "Failed to deploy project",
  DEPLOYMENT_INPROGRESS: "Failed to deploy as its already in progress",
  INTERNAL_SERVER_ERROR: "Internal server error",
} as const;

export type HttpCode = (typeof HTTP_CODES)[keyof typeof HTTP_CODES];
export type ErrorMessage = (typeof ERROR_MESSAGES)[keyof typeof ERROR_MESSAGES];
