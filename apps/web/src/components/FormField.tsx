import { Box, Typography } from "@mui/material";
import type { ReactNode } from "react";

type FormFieldProps = {
  label: string;
  htmlFor: string;
  required?: boolean;
  helperText?: string;
  children: ReactNode;
};

export const FormField = ({
  label,
  htmlFor,
  required = false,
  helperText,
  children,
}: FormFieldProps) => {
  return (
    <Box className="FormField" sx={{ width: "100%" }}>
      <Box sx={{ display: "flex", alignItems: "center", gap: 1, mb: 0.5 }}>
        <Typography
          component="label"
          htmlFor={htmlFor}
          variant="body2"
          sx={{ fontWeight: 600 }}
        >
          {label}
        </Typography>
        {required && (
          <Typography
            component="span"
            variant="caption"
            sx={{
              color: "error.main",
              border: 1,
              borderColor: "error.main",
              borderRadius: 1,
              px: 0.5,
              lineHeight: 1.6,
            }}
          >
            必須
          </Typography>
        )}
      </Box>
      {children}
      {helperText && (
        <Typography variant="caption" color="text.secondary">
          {helperText}
        </Typography>
      )}
    </Box>
  );
};
