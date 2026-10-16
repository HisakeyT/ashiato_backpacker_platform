package model

import (
	"testing"
	"time"
)

func TestTravel_Copy(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)
	planEnd := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)

	newSource := func() *Travel {
		return &Travel{
			ID:        10,
			UserID:    1,
			Title:     "京都旅行",
			IsPublic:  true,
			StartDate: start,
			EndDate:   end,
			Plans: []TravelPlan{
				{ID: 100, TravelID: 10, StartDate: start, EndDate: &planEnd, Place: "清水寺", Content: "朝イチ", SortOrder: 1},
				{ID: 101, TravelID: 10, StartDate: end, EndDate: nil, Place: "伏見稲荷", Content: "夕方", SortOrder: 2},
			},
		}
	}

	t.Run("旅行のヘッダーが複製される", func(t *testing.T) {
		src := newSource()
		got := src.Copy(2)

		if got.ID != 0 {
			t.Errorf("ID = %d, want 0", got.ID)
		}
		if got.UserID != 2 {
			t.Errorf("UserID = %d, want 2", got.UserID)
		}
		if got.IsPublic {
			t.Error("IsPublic = true, want false")
		}
		if got.Title != src.Title ||
			!got.StartDate.Equal(src.StartDate) ||
			!got.EndDate.Equal(src.EndDate) {
			t.Errorf("header not copied: %+v", got)
		}
	})

	t.Run("コピー元IDが記録され、元のIDとアドレスを共有しない", func(t *testing.T) {
		src := newSource()
		got := src.Copy(2)

		if got.CopiedFromTravelID == nil || *got.CopiedFromTravelID != 10 {
			t.Fatalf("CopiedFromTravelID = %v, want 10", got.CopiedFromTravelID)
		}
		if got.CopiedFromTravelID == &src.ID {
			t.Error("CopiedFromTravelID aliases the source ID")
		}
	})

	t.Run("プランが複製され、IDとTravelIDはゼロ値になる", func(t *testing.T) {
		src := newSource()
		got := src.Copy(2)

		if len(got.Plans) != 2 {
			t.Fatalf("len(Plans) = %d, want 2", len(got.Plans))
		}
		for i, p := range got.Plans {
			if p.ID != 0 || p.TravelID != 0 {
				t.Errorf("Plans[%d] ID/TravelID = %d/%d, want 0/0", i, p.ID, p.TravelID)
			}
			s := src.Plans[i]
			if p.Place != s.Place || p.Content != s.Content || p.SortOrder != s.SortOrder ||
				p.Location != s.Location || !p.StartDate.Equal(s.StartDate) {
				t.Errorf("Plans[%d] fields differ: got %+v, src %+v", i, p, s)
			}
		}
	})

	t.Run("プランのEndDateは深いコピーで、nilはnilのまま", func(t *testing.T) {
		src := newSource()
		got := src.Copy(2)

		if got.Plans[0].EndDate == nil || !got.Plans[0].EndDate.Equal(planEnd) {
			t.Fatalf("Plans[0].EndDate = %v, want %v", got.Plans[0].EndDate, planEnd)
		}
		if got.Plans[0].EndDate == src.Plans[0].EndDate {
			t.Error("Plans[0].EndDate pointer is shared")
		}
		if got.Plans[1].EndDate != nil {
			t.Errorf("Plans[1].EndDate = %v, want nil", got.Plans[1].EndDate)
		}
	})

	t.Run("プランが空でもパニックしない", func(t *testing.T) {
		src := &Travel{ID: 10, UserID: 1}
		got := src.Copy(2)
		if len(got.Plans) != 0 {
			t.Errorf("len(Plans) = %d, want 0", len(got.Plans))
		}
	})

	t.Run("元のデータは変更されない", func(t *testing.T) {
		src := newSource()
		_ = src.Copy(2)

		if src.ID != 10 || src.UserID != 1 || !src.IsPublic ||
			src.Plans[0].ID != 100 || src.Plans[0].TravelID != 10 {
			t.Errorf("source was mutated: %+v", src)
		}
	})
}
