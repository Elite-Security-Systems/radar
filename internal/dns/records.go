package dns

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

// RecordTypeToString converts a DNS record type code to its string representation
func RecordTypeToString(typeCode uint16) string {
	if name, exists := GetRecordTypeMapping()[typeCode]; exists {
		return name
	}
	return fmt.Sprintf("TYPE%d", typeCode)
}

// GetRecordTypeMapping returns a mapping from DNS record type numbers to their names
func GetRecordTypeMapping() map[uint16]string {
	return map[uint16]string{
		1:     "A",
		2:     "NS",
		3:     "MD",
		4:     "MF",
		5:     "CNAME",
		6:     "SOA",
		7:     "MB",
		8:     "MG",
		9:     "MR",
		10:    "NULL",
		11:    "WKS",
		12:    "PTR",
		13:    "HINFO",
		14:    "MINFO",
		15:    "MX",
		16:    "TXT",
		17:    "RP",
		18:    "AFSDB",
		19:    "X25",
		20:    "ISDN",
		21:    "RT",
		22:    "NSAP",
		23:    "NSAPPTR",
		24:    "SIG",
		25:    "KEY",
		26:    "PX",
		27:    "GPOS",
		28:    "AAAA",
		29:    "LOC",
		30:    "NXT",
		31:    "EID",
		32:    "NIMLOC",
		33:    "SRV",
		34:    "ATMA",
		35:    "NAPTR",
		36:    "KX",
		37:    "CERT",
		39:    "DNAME",
		41:    "OPT",
		42:    "APL",
		43:    "DS",
		44:    "SSHFP",
		46:    "RRSIG",
		47:    "NSEC",
		48:    "DNSKEY",
		49:    "DHCID",
		50:    "NSEC3",
		51:    "NSEC3PARAM",
		52:    "TLSA",
		53:    "SMIMEA",
		55:    "HIP",
		56:    "NINFO",
		57:    "RKEY",
		58:    "TALINK",
		59:    "CDS",
		60:    "CDNSKEY",
		61:    "OPENPGPKEY",
		62:    "CSYNC",
		99:    "SPF",
		100:   "UINFO",
		101:   "UID",
		102:   "GID",
		103:   "UNSPEC",
		104:   "NID",
		105:   "L32",
		106:   "L64",
		107:   "LP",
		108:   "EUI48",
		109:   "EUI64",
		256:   "URI",
		257:   "CAA",
		258:   "AVC",
		259:   "DOA",
		260:   "AMTRELAY",
		32768: "TA",
		32769: "DLV",
	}
}

// ExtractValue extracts the value from a DNS record based on its type
func ExtractValue(rr dns.RR) string {
	switch rr := rr.(type) {
	case *dns.A:
		return rr.A.String()
	case *dns.AAAA:
		return rr.AAAA.String()
	case *dns.CNAME:
		return rr.Target
	case *dns.MX:
		return fmt.Sprintf("%d %s", rr.Preference, rr.Mx)
	case *dns.NS:
		return rr.Ns
	case *dns.PTR:
		return rr.Ptr
	case *dns.SOA:
		return fmt.Sprintf("%s %s %d %d %d %d %d", rr.Ns, rr.Mbox, rr.Serial, rr.Refresh, rr.Retry, rr.Expire, rr.Minttl)
	case *dns.SRV:
		return fmt.Sprintf("%d %d %d %s", rr.Priority, rr.Weight, rr.Port, rr.Target)
	case *dns.TXT:
		return strings.Join(rr.Txt, " ")
	case *dns.CAA:
		return fmt.Sprintf("%d %s \"%s\"", rr.Flag, rr.Tag, rr.Value)
	case *dns.DNSKEY:
		return fmt.Sprintf("%d %d %d %s", rr.Flags, rr.Protocol, rr.Algorithm, rr.PublicKey)
	case *dns.DS:
		return fmt.Sprintf("%d %d %d %s", rr.KeyTag, rr.Algorithm, rr.DigestType, rr.Digest)
	case *dns.NAPTR:
		return fmt.Sprintf("%d %d \"%s\" \"%s\" \"%s\" %s", rr.Order, rr.Preference, rr.Flags, rr.Service, rr.Regexp, rr.Replacement)
	case *dns.RRSIG:
		return fmt.Sprintf("%s %d %d %d %d %d %d %s %s", dns.TypeToString[rr.TypeCovered], rr.Algorithm, rr.Labels, rr.OrigTtl, rr.Expiration, rr.Inception, rr.KeyTag, rr.SignerName, rr.Signature)
	case *dns.NSEC:
		return fmt.Sprintf("%s %s", rr.NextDomain, typesToString(rr.TypeBitMap))
	case *dns.TLSA:
		return fmt.Sprintf("%d %d %d %s", rr.Usage, rr.Selector, rr.MatchingType, rr.Certificate)
	default:
		return rr.String()
	}
}

// typesToString converts a slice of record types to a string
func typesToString(types []uint16) string {
	var strs []string
	for _, t := range types {
		strs = append(strs, dns.TypeToString[t])
	}
	return strings.Join(strs, " ")
}
