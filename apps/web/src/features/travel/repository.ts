import type {
  Travel,
  PublicTravel,
  CreateTravelInput,
  UpdateTravelInput,
  TravelDetail,
} from "./types";
import { apiFetch } from "../../lib/apiFetch";

export const getTravels = async (token: string): Promise<Travel[]> => {
  return apiFetch<Travel[]>("/api/travels", { token, method: "GET" });
};

export const getTravel = async (token: string, travelId: number): Promise<TravelDetail> => {
  return apiFetch<TravelDetail>(`/api/travels/${travelId}`, { token, method: "GET" });
};

export const createTravel = async (token: string, input: CreateTravelInput): Promise<Travel> => {
  return apiFetch<Travel>("/api/travels", {
    method: "POST",
    token,
    body: input,
  });
};

export const updateTravel = async (
  token: string,
  id: number,
  input: UpdateTravelInput,
): Promise<Travel> => {
  return apiFetch<Travel>(`/api/travels/${id}`, {
    method: "PATCH",
    token,
    body: input,
  });
};

export const deleteTravel = async (token: string, id: number): Promise<void> => {
  await apiFetch<void>(`/api/travels/${id}`, {
    method: "DELETE",
    token,
  });
};

export const getPublicTravels = async (): Promise<PublicTravel[]> => {
  return apiFetch<PublicTravel[]>("/api/travels/public", { method: "GET" });
};

export const getPublicTravel = async (travelId: number): Promise<PublicTravel> => {
  return apiFetch<PublicTravel>(`/api/travels/public/${travelId}`, { method: "GET" });
};

export const copyTravel = async (token: string, travelId: number): Promise<{ id: number }> => {
  return apiFetch<{ id: number }>(`/api/travels/${travelId}/copy`, {
    method: "POST",
    token,
  });
};
