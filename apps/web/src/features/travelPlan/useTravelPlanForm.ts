import { useState } from "react";
import { searchPlaces, type PlaceCandidate } from "./searchPlaces";
import type { TravelPlanFormValue } from "./TravelPlanForm";

export const usePlanDetails = (initialValue?: TravelPlanFormValue, defaultDate = "") => {
  const [startDate, setStartDate] = useState(initialValue?.startDate ?? defaultDate);
  const [endDate, setEndDate] = useState(initialValue?.endDate ?? "");
  const [place, setPlace] = useState(initialValue?.place ?? "");
  const [country, setCountry] = useState(initialValue?.country ?? "");
  const [state, setState] = useState(initialValue?.state ?? "");
  const [city, setCity] = useState(initialValue?.city ?? "");
  const [latitude, setLatitude] = useState<number | null>(initialValue?.latitude ?? null);
  const [longitude, setLongitude] = useState<number | null>(initialValue?.longitude ?? null);
  const [content, setContent] = useState(initialValue?.content ?? "");
  return {
    startDate,
    setStartDate,
    endDate,
    setEndDate,
    place,
    setPlace,
    country,
    setCountry,
    state,
    setState,
    city,
    setCity,
    latitude,
    setLatitude,
    longitude,
    setLongitude,
    content,
    setContent,
  };
};

export type PlanDetails = ReturnType<typeof usePlanDetails>;

export const usePlaceSearch = (details: PlanDetails) => {
  const [candidates, setCandidates] = useState<PlaceCandidate[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [searchError, setSearchError] = useState<string | null>(null);
  const [searched, setSearched] = useState(false);
  const handleSearch = async () => {
    if (isSearching || !details.place.trim()) return;
    setIsSearching(true);
    setSearchError(null);
    try {
      setCandidates(await searchPlaces(details.place));
      setSearched(true);
    } catch {
      setSearchError("場所の検索に失敗しました");
    } finally {
      setIsSearching(false);
    }
  };
  const handleSelect = (candidate: PlaceCandidate) => {
    details.setPlace(candidate.name);
    details.setCountry(candidate.country ?? "");
    details.setState(candidate.state ?? "");
    details.setCity(candidate.city ?? "");
    details.setLatitude(candidate.latitude);
    details.setLongitude(candidate.longitude);
    setCandidates([]);
    setSearched(false);
  };
  return { candidates, isSearching, searchError, searched, handleSearch, handleSelect };
};

export type PlaceSearchState = ReturnType<typeof usePlaceSearch>;
