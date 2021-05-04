package command

import (
	"testing"
)

const (
	strForBench = "" +
		"MFVFLVLLPLVSSQCVNLTTRTQLPPAYTNSFTRGVYYPDKVFRSSVLHSTQDLFLPFFSNVTWFHAIHVSGTNGTKRFD" +
		"NPVLPFNDGVYFASTEKSNIIRGWIFGTTLDSKTQSLLIVNNATNVVIKVCEFQFCNDPFLGVYYHKNNKSWMESEFRVY" +
		"SSANNCTFEYVSQPFLMDLEGKQGNFKNLREFVFKNIDGYFKIYSKHTPINLVRDLPQGFSALEPLVDLPIGINITRFQT" +
		"LLALHRSYLTPGDSSSGWTAGAAAYYVGYLQPRTFLLKYNENGTITDAVDCALDPLSETKCTLKSFTVEKGIYQTSNFRV" +
		"QPTESIVRFPNITNLCPFGEVFNATRFASVYAWNRKRISNCVADYSVLYNSASFSTFKCYGVSPTKLNDLCFTNVYADSF" +
		"VIRGDEVRQIAPGQTGKIADYNYKLPDDFTGCVIAWNSNNLDSKVGGNYNYLYRLFRKSNLKPFERDISTEIYQAGSTPC" +
		"NGVEGFNCYFPLQSYGFQPTNGVGYQPYRVVVLSFELLHAPATVCGPKKSTNLVKNKCVNFNFNGLTGTGVLTESNKKFL" +
		"PFQQFGRDIADTTDAVRDPQTLEILDITPCSFGGVSVITPGTNTSNQVAVLYQDVNCTEVPVAIHADQLTPTWRVYSTGS" +
		"NVFQTRAGCLIGAEHVNNSYECDIPIGAGICASYQTQTNSPRRARSVASQSIIAYTMSLGAENSVAYSNNSIAIPTNFTI" +
		"SVTTEILPVSMTKTSVDCTMYICGDSTECSNLLLQYGSFCTQLNRALTGIAVEQDKNTQEVFAQVKQIYKTPPIKDFGGF" +
		"NFSQILPDPSKPSKRSFIEDLLFNKVTLADAGFIKQYGDCLGDIAARDLICAQKFNGLTVLPPLLTDEMIAQYTSALLAG" +
		"TITSGWTFGAGAALQIPFAMQMAYRFNGIGVTQNVLYENQKLIANQFNSAIGKIQDSLSSTASALGKLQDVVNQNAQALN" +
		"TLVKQLSSNFGAISSVLNDILSRLDKVEAEVQIDRLITGRLQSLQTYVTQQLIRAAEIRASANLAATKMSECVLGQSKRV" +
		"DFCGKGYHLMSFPQSAPHGVVFLHVTYVPAQEKNFTTAPAICHDGKAHFPREGVFVSNGTHWFVTQRNFYEPQIITTDNT" +
		"FVSGNCDVVIGIVNNTVYDPLQPELDSFKEELDKYFKNHTSPDVDLGDISGINASVVNIQKEIDRLNEVAKNLNESLIDL" +
		"QELGKYEQYIKWPWYIWLGFIAGLIAIVMVTIMLCCMTSCCSCLKGCCSCGSCCKFDEDDSEPVLKGVKLHYT"
)

var (
	casesForTest = []struct {
		s      string
		substr string
		count  int
	}{
		{"", "", 0},
		{"A", "", 0},
		{"AG", "", 0},
		{"AGC", "", 0},
		{"AGCT", "", 0},
		{"AGCTU", "", 0},
		{"AGCTUAGCTU", "AGCTU", 2},
		{"AGCTUGGGGGAGCTU", "", 0},
		{"GGGGGAGCTUAGCTU", "AGCTU", 2},
		{"AGCTUAGCTUGGGGG", "AGCTU", 2},
		{"GGGGGAGCTUAGCTUGGGGG", "AGCTU", 2},
		{"AAAAAAAAAA", "", 0},
		{"あああああああ", "", 0},
		{"ひねもすのたりのたりかな", "のたり", 2},
		{"オラオラオラオラオラ", "オラ", 5},
		{"無駄無駄無駄無駄無駄", "無駄", 5},
		{"無駄あ無駄", "", 0},
		{"無駄無駄あ", "無駄", 2},
		{"無駄無駄ああ", "無駄", 2},
		{"無駄無駄あああ", "無駄", 2},
		{"無駄無駄ああああ", "無駄", 2},
		{"あ無駄無駄", "無駄", 2},
		{"ああ無駄無駄", "無駄", 2},
		{"あああ無駄無駄", "無駄", 2},
		{"ああああ無駄無駄", "無駄", 2},
		{"ああああ無駄無駄ああああ", "無駄", 2},
	}
)

func TestDetectLongestTandemRepeat1(t *testing.T) {
	for n, c := range casesForTest {
		if substr, count := detectLongestTandemRepeat1(c.s); substr != c.substr || count != c.count {
			t.Errorf("%d: %s: substr: want=%s, got=%s: count: want=%d, got=%d",
				n+1, c.s, c.substr, substr, c.count, count)
		}
	}
}

func TestDetectLongestTandemRepeat2(t *testing.T) {
	for n, c := range casesForTest {
		if substr, count := detectLongestTandemRepeat2(c.s); substr != c.substr || count != c.count {
			t.Errorf("%d: %s: substr: want=%s, got=%s: count: want=%d, got=%d",
				n+1, c.s, c.substr, substr, c.count, count)
		}
	}
}

func TestDetectLongestTandemRepeat3(t *testing.T) {
	for n, c := range casesForTest {
		if substr, count := detectLongestTandemRepeat3(c.s); substr != c.substr || count != c.count {
			t.Errorf("%d: %s: substr: want=%s, got=%s: count: want=%d, got=%d",
				n+1, c.s, c.substr, substr, c.count, count)
		}
	}
}

func BenchmarkDetectLongestTandemRepeat1(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detectLongestTandemRepeat1(strForBench)
	}
}

func BenchmarkDetectLongestTandemRepeat2(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detectLongestTandemRepeat2(strForBench)
	}
}

func BenchmarkDetectLongestTandemRepeat3(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detectLongestTandemRepeat3(strForBench)
	}
}
