export type Travel = {
  id: number;
  userId: number;
  title: string;
  startDate: string;
  endDate: string;
  isPublic: boolean;
  createdAt: string;
  updatedAt: string;
};

export type CopiedFrom = {
  travelId: number;
  authorNickname: string;
};

export type TravelDetail = Travel & {
  copiedFrom: CopiedFrom | null;
};

export interface PublicTravel extends Omit<Travel, "userId" | "isPublic"> {
  authorNickname: string;
}

export type CreateTravelInput = {
  title: string;
  startDate: string; // YYYY-MM-DD
  endDate: string;
  isPublic: boolean;
};

export type UpdateTravelInput = CreateTravelInput;
