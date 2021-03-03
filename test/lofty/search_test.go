package lofty

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"lofty/cmd/lofty"
	"testing"
)

// Test data
const shortTargetProse = "The lazy fox jumped over the fence"
const longTargetProse = "Programmers are often encouraged to use long variable names " +
	"regardless of context.\nThat is a mistake: clarity is often achieved through brevity. ∆√∫"

/**
 * Init function
 */
func TestSearch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Search")
}

/**
 * BDD Tests
 */
var _ = Describe("Search a few lines", func() {
	Describe("singular expression (A)", func() {
		Context("in line with A", func() {
			It("should be true", func() {
				Expect(lofty.SearchString("The", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("jumped", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("fence", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("Programmers", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("a", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("variable", longTargetProse)).To(Equal(true))
			})
		})

		Context("in line without A", func() {
			It("should be false", func() {
				Expect(lofty.SearchString("Thee", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("jumper", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("felt", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("Programmed", longTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("alt", longTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("variety", longTargetProse)).To(Equal(false))
			})
		})
	})
	Describe("expression with or operator (A+B)", func() {
		Context("in line with A and B", func() {
			It("should be true", func() {
				Expect(lofty.SearchString("lazy | fox", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("the | fence", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("over | the", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("are | brevity", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("is | is", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("encouraged | regardless", longTargetProse)).To(Equal(true))
			})
		})

		Context("in line with A and not B", func() {
			It("should be true", func() {
				Expect(lofty.SearchString("lazy | abc", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("the | fend", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("over | they", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("are | bravado", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("is | si", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("encouraged | regarding", longTargetProse)).To(Equal(true))
			})
		})

		Context("in line with B and not A", func() {
			It("should be true", func() {
				Expect(lofty.SearchString("laser | fox", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("than | fence", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("overtly | the", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("arsenal | brevity", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("island | is", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("encouraging | regardless", longTargetProse)).To(Equal(true))
			})
		})

		Context("in line with not B and not A", func() {
			It("should be false", func() {
				Expect(lofty.SearchString("laser | abc", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("than | fend", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("overtly | they", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("arsenal | brevado", longTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("island | si", longTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("encouraging | regarding", longTargetProse)).To(Equal(false))
			})
		})
	})
	Describe("expression with and operator (A.B) ", func() {
		Context("in line with A and B", func() {
			It("should be true", func() {
				Expect(lofty.SearchString("lazy & fox", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("The & over", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("jumped & fence", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("Programmers & mistake", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("∆√∫ & brevity", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("mistake & clarity", longTargetProse)).To(Equal(true))
			})
		})
		Context("in line with A and not B", func() {
			It("should be false", func() {
				Expect(lofty.SearchString("lazy & foxy", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("The & overt", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("jumped & fencing", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("Programmers & mistook", longTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("∆√∫ & brave", longTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("mistake & claire", longTargetProse)).To(Equal(false))
			})
		})
		Context("in line with not A and B", func() {
			It("should be false", func() {
				Expect(lofty.SearchString("lazy & foxy", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("The & overt", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("jumped & fencing", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("Programmers & mistook", longTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("through & brave", longTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("mistake & claire", longTargetProse)).To(Equal(false))
			})
		})
		Context("in line with not A and not B", func() {
			It("should be false", func() {
				Expect(lofty.SearchString("ladle & foxy", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("Thorn & overt", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("jam & fencing", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("Frog & mistook", longTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("Clark & brave", longTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("Park & claire", longTargetProse)).To(Equal(false))
			})
		})
	})

	Describe("expression with and operator A+(B+C) ", func() {
		Context("in line with A and B and C", func() {
			It("should be false", func() {
				Expect(lofty.SearchString("lazy|(the|jumped)", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("fence|(fence|fence)", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("fox|(fence|The)", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("Programmers|(mistake|variable)", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("clarity|(often|achieved)", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("use|(long|are)", longTargetProse)).To(Equal(true))
			})
		})

		Context("in line with A and B and not C", func() {
			It("should be false", func() {
				Expect(lofty.SearchString("lazy|(the|zumped)", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("fence|(fence|zence)", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("fox|(fence|zhe)", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("Programmers|(mistake|zariable)", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("∆√∫|(often|zchieved)", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("use|(long|zre)", longTargetProse)).To(Equal(true))
			})
		})

		Context("in line with A and not B and C", func() {
			It("should be false", func() {
				Expect(lofty.SearchString("lazy|(th¬kje|jumped)", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("fence|(fkjence|fence)", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("fox|(fenjce|the)", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("Programmers|(mis˚take|variable)", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("clarity|(ofte¬n|achieved)", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("use|(loßng|are)", longTargetProse)).To(Equal(true))
			})
		})

		Context("in line with not A and B and C", func() {
			It("should be false", func() {
				Expect(lofty.SearchString("ladπle|(the|jumped)", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("deønse|(fence|fence)", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("bo´´x|(fence|the)", shortTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("Grammers|(mistake|variable)", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("Cl®arity|(often|achieved)", longTargetProse)).To(Equal(true))
				Expect(lofty.SearchString("US†A|(long|are)", longTargetProse)).To(Equal(true))
			})
		})

		Context("in line with not A and not B and not C", func() {
			It("should be false", func() {
				Expect(lofty.SearchString("lazys|(these|jumpeds)", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("fencåes|(fe∂nces|fßences)", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("foxs|(fences|thes)", shortTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("Proƒgrammering|(mßistakes|variabl∂es)", longTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("claritys|(oftens|achi∫eveds)", longTargetProse)).To(Equal(false))
				Expect(lofty.SearchString("uses|(longs|ares)", longTargetProse)).To(Equal(false))
			})
		})

		Describe("expression with and operator (A+B)+C ", func() {
			Context("in line with A and B and C", func() {
				It("should be false", func() {
					Expect(lofty.SearchString("(the|jumped)|lazy", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(fence|fence)|fence", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(fence|The)|fox", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(mistake|variable)|Programmers", longTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(often|achieved)|clarity", longTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(long|are)|use", longTargetProse)).To(Equal(true))
				})
			})

			Context("in line with A and B and not C", func() {
				It("should be false", func() {
					Expect(lofty.SearchString("(the|zumped)|lazy", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(fence|zence)|fence", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(fence|zhe)|fox", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(mistake|zariable)|Programmers", longTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(often|zchieved)|∆√∫", longTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(long|zre)|use", longTargetProse)).To(Equal(true))
				})
			})

			Context("in line with A and not B and C", func() {
				It("should be false", func() {
					Expect(lofty.SearchString("(th¬kje|jumped)|lazy", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(fkjence|fence)|fence", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(fenjce|the)|fox", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(mis˚take|variable)|Programmers", longTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(ofte¬n|achieved)|clarity", longTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(loßng|are)|use", longTargetProse)).To(Equal(true))
				})
			})

			Context("in line with not A and B and C", func() {
				It("should be false", func() {
					Expect(lofty.SearchString("(the|jumped)|ladπle", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(fence|fence)|deønse", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(fence|the)|bo´´x", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(mistake|variable)|Grammers", longTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(often|achieved)|Cl®arity", longTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(long|are)|US†A", longTargetProse)).To(Equal(true))
				})
			})

			Context("in line with not A and not B and not C", func() {
				It("should be false", func() {
					Expect(lofty.SearchString("(these|jumpeds)|lazys", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(fe∂nces|fßences)|fencåes", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(fences|thes)|foxs", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(mßistakes|variabl∂es)|Proƒgrammering", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(oftens|achi∫eveds)|claritys", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(longs|ares)|uses", longTargetProse)).To(Equal(false))
				})
			})
		})


		Describe("expression with and operator (A.B).C ", func() {
			Context("in line with A and B and C", func() {
				It("should be false", func() {
					Expect(lofty.SearchString("(the&jumped)&lazy", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(fence&fence)&fence", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(fence&The)&fox", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(mistake&variable)&Programmers", longTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(often&achieved)&clarity", longTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("(long&are)&use", longTargetProse)).To(Equal(true))
				})
			})

			Context("in line with A and B and not C", func() {
				It("should be false", func() {
					Expect(lofty.SearchString("(the&zumped)&lazy", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(fence&zence)&fence", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(fence&zhe)&fox", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(mistake&zariable)&Programmers", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(often&zchieved)&∆√∫", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(long&zre)&use", longTargetProse)).To(Equal(false))
				})
			})

			Context("in line with A and not B and C", func() {
				It("should be false", func() {
					Expect(lofty.SearchString("(th¬kje&jumped)&lazy", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(fkjence&fence)&fence", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(fenjce&the)&fox", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(mis˚take&variable)&Programmers", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(ofte¬n&achieved)&clarity", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(loßng&are)&use", longTargetProse)).To(Equal(false))
				})
			})

			Context("in line with not A and B and C", func() {
				It("should be false", func() {
					Expect(lofty.SearchString("(the&jumped)&ladπle", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(fence&fence)&deønse", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(fence&the)&bo´´x", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(mistake&variable)&Grammers", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(often&achieved)&Cl®arity", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(long&are)&US†A", longTargetProse)).To(Equal(false))
				})
			})

			Context("in line with not A and not B and not C", func() {
				It("should be false", func() {
					Expect(lofty.SearchString("(these&jumpeds)&lazys", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(fe∂nces&fßences)&fencåes", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(fences&thes)&foxs", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(mßistakes&variabl∂es)&Proƒgrammering", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(oftens&achi∫eveds)&claritys", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("(longs&ares)&uses", longTargetProse)).To(Equal(false))
				})
			})
		})

		Describe("expression with and operator ((A.B).C)", func() {
			Context("in line with A and B and C)", func() {
				It("should be false)", func() {
					Expect(lofty.SearchString("((fence&fence)&fence)", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("((the&jumped)&lazy)", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("((fence&The)&fox)", shortTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("((mistake&variable)&Programmers)", longTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("((often&achieved)&clarity)", longTargetProse)).To(Equal(true))
					Expect(lofty.SearchString("((long&are)&use)", longTargetProse)).To(Equal(true))
				})
			})

			Context("(in line with A and B and not C)", func() {
				It("(should be false)", func() {
					Expect(lofty.SearchString("((the&zumped)&lazy)", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((fence&zence)&fence)", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((fence&zhe)&fox)", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((mistake&zariable)&Programmers)", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((often&zchieved)&∆√∫)", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((long&zre)&use)", longTargetProse)).To(Equal(false))
				})
			})

			Context("(in line with A and not B and C)", func() {
				It("(should be false)", func() {
					Expect(lofty.SearchString("((th¬kje&jumped)&lazy)", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((fkjence&fence)&fence)", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((fenjce&the)&fox)", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((mis˚take&variable)&Programmers)", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((ofte¬n&achieved)&clarity)", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((loßng&are)&use)", longTargetProse)).To(Equal(false))
				})
			})

			Context("(in line with not A and B and C)", func() {
				It("(should be false)", func() {
					Expect(lofty.SearchString("((the&jumped)&ladπle)", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((fence&fence)&deønse)", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((fence&the)&bo´´x)", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((mistake&variable)&Grammers)", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((often&achieved)&Cl®arity)", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((long&are)&US†A)", longTargetProse)).To(Equal(false))
				})
			})

			Context("(in line with not A and not B and not C)", func() {
				It("(should be false)", func() {
					Expect(lofty.SearchString("((these&jumpeds)&lazys)", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((fe∂nces&fßences)&fencåes)", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((fences&thes)&foxs)", shortTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((mßistakes&variabl∂es)&Proƒgrammering)", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((oftens&achi∫eveds)&claritys)", longTargetProse)).To(Equal(false))
					Expect(lofty.SearchString("((longs&ares)&uses)", longTargetProse)).To(Equal(false))
				})
			})
		})
	})
})
