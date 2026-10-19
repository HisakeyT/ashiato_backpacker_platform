export type TravelPlan = {
  id: number;
  travelId: number;
  startDate: string | null;
  endDate: string | null; // NULL は startDate と同日
  place: string;
  country: string | null;
  state: string | null;
  city: string | null;
  latitude: number | null;
  longitude: number | null;
  content: string;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
};

export type CreateTravelPlanInput = {
  startDate: string | null; // YYYY-MM-DD or undetermined
  endDate?: string | null; // YYYY-MM-DD。省略可(省略 = 同日)
  place: string;
  country?: string;
  state?: string;
  city?: string;
  latitude?: number;
  longitude?: number;
  content: string;
  sortOrder: number;
};

export type UpdateTravelPlanInput = {
  startDate?: string | null; // YYYY-MM-DD
  endDate?: string | null; // YYYY-MM-DD
  place?: string;
  country?: string;
  state?: string;
  city?: string;
  latitude?: number;
  longitude?: number;
  content?: string;
};
