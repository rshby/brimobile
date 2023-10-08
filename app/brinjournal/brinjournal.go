package brinjournal

type BrinJournal struct {
	Id         int    `json:"id,omitempty"`
	BranchCode string `json:"branch_code,omitempty"`
	JournalSeq int    `json:"journal_seq,omitempty"`
}
