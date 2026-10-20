import camelcaseKeys from "camelcase-keys";
import snakecaseKeys from "snakecase-keys";

type HTTPMethod = "GET" | "POST" | "PUT" | "DELETE" | "PATCH";

type ApiFetchOptions = {
  method: HTTPMethod;
  token?: string;
  body?: Record<string, unknown>;
};

export class ApiError extends Error {
  status: number;
  method: string;
  path: string;

  constructor(method: string, path: string, status: number) {
    super(`Request failed: ${method} ${path} (${status})`);
    this.name = "ApiError";
    this.status = status;
    this.method = method;
    this.path = path;
  }
}

const createHeaders = ({ token, body }: ApiFetchOptions): Record<string, string> => {
  const headers: Record<string, string> = {};
  if (token) headers.Authorization = `Bearer ${token}`;
  if (body) headers["Content-Type"] = "application/json";
  return headers;
};

export const apiFetch = async <T>(path: string, options: ApiFetchOptions): Promise<T> => {
  const { method, body } = options;

  const response = await fetch(path, {
    method,
    headers: createHeaders(options),
    body: body ? JSON.stringify(snakecaseKeys(body, { deep: true })) : undefined,
  });

  if (!response.ok) {
    if (response.status === 401 && options.token) {
      window.dispatchEvent(new Event("auth:unauthorized"));
    }
    throw new ApiError(method, path, response.status);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  const data = await response.json();
  return camelcaseKeys(data, { deep: true }) as T;
};
