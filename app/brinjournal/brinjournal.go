package brinjournal

type BrinJournal struct {
	Id         int    `json:"id,omitempty"`
	BranchCode string `json:"branch_code,omitempty"`
	JournalSeq string `json:"journal_seq,omitempty"`
}
