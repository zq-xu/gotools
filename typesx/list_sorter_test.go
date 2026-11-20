package typesx

import "testing"

func TestNewSorterValidQuery(t *testing.T) {
	s := NewSorter("name,asc:::created_at,desc")

	if got, want := s.MysqlString(), "name asc,created_at desc"; got != want {
		t.Fatalf("MysqlString() = %q, want %q", got, want)
	}
}

func TestNewSorterDropsInjectedCondition(t *testing.T) {
	cases := []string{
		"(select sleep(5))",
		"name;drop table accounts",
		"1=1",
		"`id`",
		"",
	}

	for _, in := range cases {
		if got, want := NewSorter(in).MysqlString(), "updated_at desc"; got != want {
			t.Errorf("NewSorter(%q).MysqlString() = %q, want %q (should fall back to default)", in, got, want)
		}
	}
}

func TestNewSorterKeepsValidPartOfMixedQuery(t *testing.T) {
	// 第二个排序键非法，只保留第一个
	s := NewSorter("name,asc:::(select sleep(5))")

	if got, want := s.MysqlString(), "name asc"; got != want {
		t.Fatalf("MysqlString() = %q, want %q", got, want)
	}
}

func TestNewSorterNormalizesOrder(t *testing.T) {
	s := NewSorter("name,ascending:::age,descending:::id,nonsense")

	if got, want := s.MysqlString(), "name asc,age desc,id asc"; got != want {
		t.Fatalf("MysqlString() = %q, want %q", got, want)
	}
}
