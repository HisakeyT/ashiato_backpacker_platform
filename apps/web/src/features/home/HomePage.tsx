import EditNoteIcon from "@mui/icons-material/EditNote";
import ExploreIcon from "@mui/icons-material/Explore";
import PublicIcon from "@mui/icons-material/Public";
import { Box, Button, Card, CardContent, Container, Typography } from "@mui/material";
import { Link as RouterLink } from "react-router-dom";

import { useAuth } from "../../app/providers/useAuth";
import ctaImage from "../../assets/ctaImage.webp";
import heroImage from "../../assets/hero.webp";
import { SERVICE_NAME } from "../../constants/site.ts";

type Cta = { label: string; to: string };

const features = [
  {
    title: "日ごとに記録する",
    description: "行った場所とメモを日付ごとに整理。長い旅でも、あとから振り返りやすくなります。",
    icon: <EditNoteIcon fontSize="large" color="primary" />,
  },
  {
    title: "公開して共有する",
    description: "旅行ごとに公開・非公開を選べます。非公開なら、自分用の旅ノートとして使えます。",
    icon: <PublicIcon fontSize="large" color="primary" />,
  },
  {
    title: "先人の旅から学ぶ",
    description: "公開された旅行記を見て、ルートや滞在日数の参考にできます。",
    icon: <ExploreIcon fontSize="large" color="primary" />,
  },
];

const heroSx = {
  position: "relative",
  color: "common.white",
  py: { xs: 10, md: 16 },
  backgroundImage: `url(${heroImage})`,
  backgroundSize: "cover",
  backgroundPosition: { xs: "75% center", md: "center" },
  "&::before": {
    content: '""',
    position: "absolute",
    inset: 0,
    background: {
      xs: "rgba(10,25,50,0.6)",
      md: "linear-gradient(90deg, rgba(10,25,50,0.65) 0%, rgba(10,25,50,0.45) 45%, rgba(10,25,50,0.1) 75%)",
    },
  },
} as const;

const ctaSx = {
  position: "relative",
  py: { xs: 8, md: 12 },
  color: "common.white",
  backgroundImage: `url(${ctaImage})`,
  backgroundSize: "cover",
  backgroundPosition: { xs: "80% center", md: "center 50%" },
  "&::before": {
    content: '""',
    position: "absolute",
    inset: 0,
    background:
      "linear-gradient(90deg, rgba(20,45,85,0.8) 0%, rgba(30,58,95,0.55) 50%, rgba(56,116,203,0.15) 100%)",
  },
} as const;

const HeroIntro = () => (
  <>
    <Typography
      variant="h2"
      component="h1"
      sx={{
        fontSize: { xs: "2rem", md: "3.25rem" },
        fontWeight: 700,
        wordBreak: "keep-all",
        overflowWrap: "anywhere",
        mb: 2,
      }}
    >
      次のバックパッカーの、地図になる。
    </Typography>
    <Typography
      sx={{
        fontSize: { xs: "1rem", md: "1.25rem" },
        mb: 4,
        maxWidth: 560,
        wordBreak: "keep-all",
        overflowWrap: "anywhere",
      }}
    >
      {SERVICE_NAME} は、バックパッカーのための旅の記録・共有サービスです。
      日ごとのプランを残し、公開して、これから旅に出る人の道しるべにしよう。
    </Typography>
  </>
);

const HeroActions = ({ cta }: { cta: Cta }) => (
  <Box sx={{ display: "flex", gap: 2, flexWrap: "wrap" }}>
    <Button
      component={RouterLink}
      to={cta.to}
      variant="contained"
      size="large"
      sx={{ bgcolor: "common.white", color: "primary.main", "&:hover": { bgcolor: "grey.100" } }}
    >
      {cta.label}
    </Button>
    <Button
      component={RouterLink}
      to="/public/travels"
      variant="outlined"
      size="large"
      color="inherit"
    >
      みんなの旅行記を見る
    </Button>
  </Box>
);

const HomeHero = ({ cta }: { cta: Cta }) => (
  <Box className="HomePage__hero" sx={heroSx}>
    <Container maxWidth="md" sx={{ position: "relative" }}>
      <HeroIntro />
      <HeroActions cta={cta} />
    </Container>
  </Box>
);

const HomeFeatures = () => (
  <Container maxWidth="lg" sx={{ py: { xs: 6, md: 10 } }}>
    <Typography
      variant="h4"
      component="h2"
      align="center"
      sx={{ mb: 5, fontSize: { xs: "1.5rem", md: "2.125rem" } }}
    >
      旅の記録を、もっと自由に
    </Typography>
    <Box sx={{ display: "grid", gap: 3, gridTemplateColumns: { xs: "1fr", md: "repeat(3, 1fr)" } }}>
      {features.map((feature) => (
        <Card key={feature.title} className="HomePage__featureCard" sx={{ height: "100%" }}>
          <CardContent sx={{ p: 3 }}>
            <Box sx={{ color: "primary.main", mb: 1.5, "& svg": { fontSize: 40 } }}>
              {feature.icon}
            </Box>
            <Typography variant="h6" component="h3" gutterBottom>
              {feature.title}
            </Typography>
            <Typography
              color="text.secondary"
              sx={{ wordBreak: "keep-all", overflowWrap: "anywhere" }}
            >
              {feature.description}
            </Typography>
          </CardContent>
        </Card>
      ))}
    </Box>
  </Container>
);

const HomeCta = ({ cta }: { cta: Cta }) => (
  <Box sx={ctaSx}>
    <Container maxWidth="sm" sx={{ position: "relative", textAlign: "center" }}>
      <Typography variant="h5" component="h2" sx={{ mb: 3 }}>
        次の旅の前に、誰かの旅をのぞいてみよう。
      </Typography>
      <Button
        component={RouterLink}
        to={cta.to}
        variant="contained"
        size="large"
        sx={{ bgcolor: "common.white", color: "primary.main", "&:hover": { bgcolor: "grey.100" } }}
      >
        {cta.label}
      </Button>
    </Container>
  </Box>
);

export const HomePage = () => {
  const { isAuthenticated } = useAuth();
  const cta = isAuthenticated
    ? { label: "自分の旅行へ", to: "/travels" }
    : { label: "はじめる", to: "/register" };
  return (
    <Box className="HomePage">
      <HomeHero cta={cta} />
      <HomeFeatures />
      <HomeCta cta={cta} />
    </Box>
  );
};
