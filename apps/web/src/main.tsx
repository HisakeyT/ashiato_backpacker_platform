import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider } from "react-router-dom";

import { AuthProvider } from "./app/providers/AuthProvider";
import { FlashProvider } from "./app/providers/FlashProvider";
import { router } from "./app/routes";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <AuthProvider>
      <FlashProvider>
        <RouterProvider router={router} />
      </FlashProvider>
    </AuthProvider>
  </StrictMode>,
);
