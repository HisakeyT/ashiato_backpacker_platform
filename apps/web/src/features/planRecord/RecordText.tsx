import { Link, Typography } from "@mui/material";

const safeURL = (value: string): boolean => {
  try {
    const url = new URL(value);
    return (
      (url.protocol === "https:" || url.protocol === "http:") && !url.username && !url.password
    );
  } catch {
    return false;
  }
};

export const RecordText = ({ text }: { text: string }) => (
  <Typography variant="body2" sx={{ whiteSpace: "pre-wrap", overflowWrap: "anywhere" }}>
    {text.split(/(https?:\/\/[^\s<>"']+)/g).map((part, index) =>
      index % 2 === 1 && safeURL(part) ? (
        <Link key={index} href={part} target="_blank" rel="noopener noreferrer">
          {part}
        </Link>
      ) : (
        part
      ),
    )}
  </Typography>
);
