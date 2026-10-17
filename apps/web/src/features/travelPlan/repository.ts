import { apiFetch } from "../../lib/apiFetch";
import type { CreateTravelPlanInput, TravelPlan, UpdateTravelPlanInput } from "./types";

export const getTravelPlans = async (token: string, travelId: number): Promise<TravelPlan[]> => {
  return apiFetch<TravelPlan[]>(`/api/travels/${travelId}/plans`, { token, method: "GET" });
};

export const createTravelPlan = async (
  token: string,
  travelId: number,
  input: CreateTravelPlanInput,
): Promise<TravelPlan> => {
  return apiFetch<TravelPlan>(`/api/travels/${travelId}/plans`, {
    method: "POST",
    token,
    body: input,
  });
};

export const updateTravelPlan = async (
  token: string,
  ids: { travelId: number; planId: number },
  input: UpdateTravelPlanInput,
): Promise<TravelPlan> => {
  return apiFetch<TravelPlan>(`/api/travels/${ids.travelId}/plans/${ids.planId}`, {
    method: "PATCH",
    token,
    body: input,
  });
};

export const deleteTravelPlan = async (
  token: string,
  travelId: number,
  planId: number,
): Promise<void> => {
  await apiFetch<void>(`/api/travels/${travelId}/plans/${planId}`, {
    method: "DELETE",
    token,
  });
};

export const getPublicTravelPlans = async (travelId: number): Promise<TravelPlan[]> => {
  return apiFetch<TravelPlan[]>(`/api/travels/public/${travelId}/plans`, { method: "GET" });
};
