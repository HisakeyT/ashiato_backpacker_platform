export type TravelPlan = {
  id: number;
  travelId: number;
  startDate: string;
  place: string;
  content: string;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
};

export type CreateTravelPlanInput = {
  startDate: string; // YYYY-MM-DD
  place: string;
  content: string;
  sortOrder: number;
}

export type UpdateTravelPlanInput = {
  startDate?: string; // YYYY-MM-DD
  place?: string;
  content?: string;
}
