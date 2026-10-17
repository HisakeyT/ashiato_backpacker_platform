import React from "react";
import { Breadcrumbs, Link, Typography } from "@mui/material";
import { Link as RouterLink } from "react-router-dom";

export type BreadcrumbItem = {
  label: string;
  to?: string;
};

type Props = {
  items: BreadcrumbItem[];
};

export const PageBreadcrumbs = ({ items }: Props) => {
  return (
    <Breadcrumbs aria-label="パンくずリスト" sx={{ mb: 2 }}>
      {items.map((item, index) => {
        const isLast = index === items.length - 1;

        // 最後の要素(現在のページ)と、リンク先のない要素は文字だけ
        if (isLast || !item.to) {
          return (
            <Typography
              key={`${index}-${item.label}`}
              color="text.primary"
              sx={{ wordBreak: "break-word", fontWeight: 600 }}
            >
              {item.label}
            </Typography>
          );
        }

        return (
          <Link
            key={`${index}-${item.label}`}
            component={RouterLink}
            to={item.to}
            underline="hover"
          >
            {item.label}
          </Link>
        );
      })}
    </Breadcrumbs>
  );
}
