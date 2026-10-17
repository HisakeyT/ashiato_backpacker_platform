import { createBrowserRouter } from "react-router-dom";

import { AppLayout } from "./AppLayout";
import { HomePage } from "../features/home/HomePage";
import { RegisterPage } from "../features/auth/RegisterPage";
import { LoginPage } from "../features/auth/LoginPage";
import { RequireAuth } from "./RequireAuth";
import { TravelListPage } from "../features/travel/TravelListPage";
import { TravelDetailPage } from "../features/travel/TravelDetailPage";
import { TravelCreatePage } from "../features/travel/TravelCreatePage";
import { TravelEditPage } from "../features/travel/TravelEditPage";
import { PublicTravelListPage } from "../features/publicTravel/PublicTravelListPage";
import { PublicTravelDetailPage } from "../features/publicTravel/PublicTravelDetailPage";

export const router = createBrowserRouter([
  {
    element: <AppLayout />,
    children: [
      { path: "/", element: <HomePage /> },
      { path: "/register", element: <RegisterPage /> },
      { path: "/login", element: <LoginPage /> },
      { path: "/public/travels", element: <PublicTravelListPage /> },
      { path: "/public/travels/:travelId", element: <PublicTravelDetailPage /> },
      {
        element: <RequireAuth />,
        children: [
          { path: "/travels", element: <TravelListPage /> },
          { path: "/travels/new", element: <TravelCreatePage /> },
          { path: "/travels/:travelId/edit", element: <TravelEditPage /> },
          { path: "/travels/:travelId", element: <TravelDetailPage /> },
        ],
      },
    ],
  },
]);
