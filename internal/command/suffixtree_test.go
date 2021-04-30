package command

import (
	"testing"
)

func TestDetectLongestTandemRepeat(t *testing.T) {
	cases := []struct {
		s      string
		substr string
		count  int
	}{
		{"AGCTU", "", 0},
		{"AGCTUAGCTU", "AGCTU", 2},
		{"AGCTUUTCGAAGCTU", "", 0},
		{"AGCTUUTCGAAGCTUAGCTU", "AGCTU", 2},
		{"AAAAAAAAAAAAAA", "A", 14},
		{"ああああああああ", "あ", 8},
		{"ひねもすのたりのたりかな", "のたり", 2},
		{"オラオラオラオラオラ", "オラ", 5},
		{"無駄無駄ああああ", "無駄", 2},
		{"", "", 0},
	}

	for n, c := range cases {
		if substr, count := detectLongestTandemRepeat(c.s); substr != c.substr || count != c.count {
			t.Errorf("%d: %s: substr: want=%s, got=%s: count: want=%d, got=%d", n+1, c.s, c.substr, substr, c.count, count)
		}
	}
}

func BenchmarkDetectLongestTandemRepeat(b *testing.B) {
	s := "" +
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = detectLongestTandemRepeat(s)
	}
}