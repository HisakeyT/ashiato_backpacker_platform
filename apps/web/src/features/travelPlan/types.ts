export type TravelPlan = {
  id: number;
  travelId: number;
  startDate: string;
  endDate: string | null; // NULL は startDate と同日
  place: string;
  content: string;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
};

export type CreateTravelPlanInput = {
  startDate: string; // YYYY-MM-DD
  endDate?: string;  // YYYY-MM-DD。省略可(省略 = 同日)
  place: string;
  content: string;
  sortOrder: number;
}

export type UpdateTravelPlanInput = {
  startDate?: string; // YYYY-MM-DD
  endDate?: string;   // YYYY-MM-DD
  place?: string;
  content?: string;
}
