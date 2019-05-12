package io

import (
	"log"
	"time"

	"github.com/djhworld/simple-computer/circuit"
	"github.com/djhworld/simple-computer/components"
)

const BUS_WIDTH = 16

type KeyPress struct {
	Value  int
	IsDown bool
}

// [cpu] <-------------> keyboard adapter <----------- keyboard <----------- [keyPressChannel]
//         read/write                        write                 notify
type KeyboardAdapter struct {
	KeyboardInBus *components.Bus

	ioBus   *components.IOBus
	mainBus *components.Bus

	memoryBit       *components.Bit
	keycodeRegister components.Register

	andGate1            components.ANDGate8
	notGatesForAndGate1 [4]circuit.NOTGate

	andGate2            components.ANDGate3
	andGate3            components.ANDGate3
	notGatesForAndGate3 [2]circuit.NOTGate

	andGate4 circuit.ANDGate
}

func NewKeyboardAdapter() *KeyboardAdapter {
	k := new(KeyboardAdapter)
	k.KeyboardInBus = components.NewBus(BUS_WIDTH)
	return k
}

func (k *KeyboardAdapter) Connect(ioBus *components.IOBus, mainBus *components.Bus) {
	k.ioBus = ioBus
	k.mainBus = mainBus
	k.memoryBit = components.NewBit()
	k.memoryBit.Update(false, true)
	k.memoryBit.Update(false, false)
	k.andGate1 = *components.NewANDGate8()
	k.andGate2 = *components.NewANDGate3()
	k.andGate3 = *components.NewANDGate3()
	k.andGate4 = *circuit.NewANDGate()
	k.keycodeRegister = *components.NewRegister("KCR", k.KeyboardInBus, k.mainBus)

	for i := range k.notGatesForAndGate1 {
		k.notGatesForAndGate1[i] = *circuit.NewNOTGate()
	}

	for i := range k.notGatesForAndGate3 {
		k.notGatesForAndGate3[i] = *circuit.NewNOTGate()
	}
}

func (k *KeyboardAdapter) Update() {
	k.updateKeycodeReg()
	k.update()
}

func (k *KeyboardAdapter) update() {
	k.notGatesForAndGate1[0].Update(k.mainBus.GetOutputWire(8))
	k.notGatesForAndGate1[1].Update(k.mainBus.GetOutputWire(9))
	k.notGatesForAndGate1[2].Update(k.mainBus.GetOutputWire(10))
	k.notGatesForAndGate1[3].Update(k.mainBus.GetOutputWire(11))

	k.andGate1.Update(
		k.notGatesForAndGate1[0].Output(),
		k.notGatesForAndGate1[1].Output(),
		k.notGatesForAndGate1[2].Output(),
		k.notGatesForAndGate1[3].Output(),
		k.mainBus.GetOutputWire(12),
		k.mainBus.GetOutputWire(13),
		k.mainBus.GetOutputWire(14),
		k.mainBus.GetOutputWire(15),
	)

	//TODO update these to use the helper methods
	k.andGate2.Update(
		k.ioBus.GetOutputWire(components.CLOCK_SET),
		k.ioBus.GetOutputWire(components.DATA_OR_ADDRESS),
		k.ioBus.GetOutputWire(components.MODE),
	)
	k.memoryBit.Update(k.andGate1.Output(), k.andGate2.Output())

	k.notGatesForAndGate3[0].Update(k.ioBus.GetOutputWire(components.DATA_OR_ADDRESS))
	k.notGatesForAndGate3[1].Update(k.ioBus.GetOutputWire(components.MODE))

	k.andGate3.Update(
		k.ioBus.GetOutputWire(components.CLOCK_ENABLE),
		k.notGatesForAndGate3[0].Output(),
		k.notGatesForAndGate3[1].Output(),
	)

	k.andGate4.Update(k.memoryBit.Get(), k.andGate3.Output())
}

func (k *KeyboardAdapter) updateKeycodeReg() {
	if k.andGate4.Output() {
		k.keycodeRegister.Set()

		k.keycodeRegister.Enable()
		k.keycodeRegister.Update()
		k.keycodeRegister.Disable()

		// clear the register once everything is out
		k.KeyboardInBus.SetValue(0x00)
		k.keycodeRegister.Update()
		k.keycodeRegister.Unset()
		k.keycodeRegister.Update()
	}
}

type Keyboard struct {
	outBus          *components.Bus
	keyPressChannel chan *KeyPress
	quit            chan bool
}

func NewKeyboard(keyPressChannel chan *KeyPress, quit chan bool) *Keyboard {
	k := new(Keyboard)
	k.keyPressChannel = keyPressChannel
	k.quit = quit
	return k
}

func (k *Keyboard) ConnectTo(bus *components.Bus) {
	log.Println("Connecting keyboard to bus")
	k.outBus = bus
}

func (k *Keyboard) Run() {
	clock := time.Tick(33 * time.Millisecond)
	for {
		<-clock
		select {
		case <-k.quit:
			log.Println("Stopping keyboard")
			return
		case key := <-k.keyPressChannel:
			if key.IsDown {
				k.outBus.SetValue(uint16(key.Value))
			}
		}
	}
}

var up_key_presses = []KeyPress{
	KeyPress{1, false},
	KeyPress{2, false},
	KeyPress{3, false},
	KeyPress{4, false},
	KeyPress{5, false},
	KeyPress{6, false},
	KeyPress{7, false},
	KeyPress{8, false},
	KeyPress{9, false},
	KeyPress{10, false},
	KeyPress{11, false},
	KeyPress{12, false},
	KeyPress{13, false},
	KeyPress{14, false},
	KeyPress{15, false},
	KeyPress{16, false},
	KeyPress{17, false},
	KeyPress{18, false},
	KeyPress{19, false},
	KeyPress{20, false},
	KeyPress{21, false},
	KeyPress{22, false},
	KeyPress{23, false},
	KeyPress{24, false},
	KeyPress{25, false},
	KeyPress{26, false},
	KeyPress{27, false},
	KeyPress{28, false},
	KeyPress{29, false},
	KeyPress{30, false},
	KeyPress{31, false},
	KeyPress{32, false},
	KeyPress{33, false},
	KeyPress{34, false},
	KeyPress{35, false},
	KeyPress{36, false},
	KeyPress{37, false},
	KeyPress{38, false},
	KeyPress{39, false},
	KeyPress{40, false},
	KeyPress{41, false},
	KeyPress{42, false},
	KeyPress{43, false},
	KeyPress{44, false},
	KeyPress{45, false},
	KeyPress{46, false},
	KeyPress{47, false},
	KeyPress{48, false},
	KeyPress{49, false},
	KeyPress{50, false},
	KeyPress{51, false},
	KeyPress{52, false},
	KeyPress{53, false},
	KeyPress{54, false},
	KeyPress{55, false},
	KeyPress{56, false},
	KeyPress{57, false},
	KeyPress{58, false},
	KeyPress{59, false},
	KeyPress{60, false},
	KeyPress{61, false},
	KeyPress{62, false},
	KeyPress{63, false},
	KeyPress{64, false},
	KeyPress{65, false},
	KeyPress{66, false},
	KeyPress{67, false},
	KeyPress{68, false},
	KeyPress{69, false},
	KeyPress{70, false},
	KeyPress{71, false},
	KeyPress{72, false},
	KeyPress{73, false},
	KeyPress{74, false},
	KeyPress{75, false},
	KeyPress{76, false},
	KeyPress{77, false},
	KeyPress{78, false},
	KeyPress{79, false},
	KeyPress{80, false},
	KeyPress{81, false},
	KeyPress{82, false},
	KeyPress{83, false},
	KeyPress{84, false},
	KeyPress{85, false},
	KeyPress{86, false},
	KeyPress{87, false},
	KeyPress{88, false},
	KeyPress{89, false},
	KeyPress{90, false},
	KeyPress{91, false},
	KeyPress{92, false},
	KeyPress{93, false},
	KeyPress{94, false},
	KeyPress{95, false},
	KeyPress{96, false},
	KeyPress{97, false},
	KeyPress{98, false},
	KeyPress{99, false},
	KeyPress{100, false},
	KeyPress{101, false},
	KeyPress{102, false},
	KeyPress{103, false},
	KeyPress{104, false},
	KeyPress{105, false},
	KeyPress{106, false},
	KeyPress{107, false},
	KeyPress{108, false},
	KeyPress{109, false},
	KeyPress{110, false},
	KeyPress{111, false},
	KeyPress{112, false},
	KeyPress{113, false},
	KeyPress{114, false},
	KeyPress{115, false},
	KeyPress{116, false},
	KeyPress{117, false},
	KeyPress{118, false},
	KeyPress{119, false},
	KeyPress{120, false},
	KeyPress{121, false},
	KeyPress{122, false},
	KeyPress{123, false},
	KeyPress{124, false},
	KeyPress{125, false},
	KeyPress{126, false},
	KeyPress{127, false},
	KeyPress{128, false},
	KeyPress{129, false},
	KeyPress{130, false},
	KeyPress{131, false},
	KeyPress{132, false},
	KeyPress{133, false},
	KeyPress{134, false},
	KeyPress{135, false},
	KeyPress{136, false},
	KeyPress{137, false},
	KeyPress{138, false},
	KeyPress{139, false},
	KeyPress{140, false},
	KeyPress{141, false},
	KeyPress{142, false},
	KeyPress{143, false},
	KeyPress{144, false},
	KeyPress{145, false},
	KeyPress{146, false},
	KeyPress{147, false},
	KeyPress{148, false},
	KeyPress{149, false},
	KeyPress{150, false},
	KeyPress{151, false},
	KeyPress{152, false},
	KeyPress{153, false},
	KeyPress{154, false},
	KeyPress{155, false},
	KeyPress{156, false},
	KeyPress{157, false},
	KeyPress{158, false},
	KeyPress{159, false},
	KeyPress{160, false},
	KeyPress{161, false},
	KeyPress{162, false},
	KeyPress{163, false},
	KeyPress{164, false},
	KeyPress{165, false},
	KeyPress{166, false},
	KeyPress{167, false},
	KeyPress{168, false},
	KeyPress{169, false},
	KeyPress{170, false},
	KeyPress{171, false},
	KeyPress{172, false},
	KeyPress{173, false},
	KeyPress{174, false},
	KeyPress{175, false},
	KeyPress{176, false},
	KeyPress{177, false},
	KeyPress{178, false},
	KeyPress{179, false},
	KeyPress{180, false},
	KeyPress{181, false},
	KeyPress{182, false},
	KeyPress{183, false},
	KeyPress{184, false},
	KeyPress{185, false},
	KeyPress{186, false},
	KeyPress{187, false},
	KeyPress{188, false},
	KeyPress{189, false},
	KeyPress{190, false},
	KeyPress{191, false},
	KeyPress{192, false},
	KeyPress{193, false},
	KeyPress{194, false},
	KeyPress{195, false},
	KeyPress{196, false},
	KeyPress{197, false},
	KeyPress{198, false},
	KeyPress{199, false},
	KeyPress{200, false},
	KeyPress{201, false},
	KeyPress{202, false},
	KeyPress{203, false},
	KeyPress{204, false},
	KeyPress{205, false},
	KeyPress{206, false},
	KeyPress{207, false},
	KeyPress{208, false},
	KeyPress{209, false},
	KeyPress{210, false},
	KeyPress{211, false},
	KeyPress{212, false},
	KeyPress{213, false},
	KeyPress{214, false},
	KeyPress{215, false},
	KeyPress{216, false},
	KeyPress{217, false},
	KeyPress{218, false},
	KeyPress{219, false},
	KeyPress{220, false},
	KeyPress{221, false},
	KeyPress{222, false},
	KeyPress{223, false},
	KeyPress{224, false},
	KeyPress{225, false},
	KeyPress{226, false},
	KeyPress{227, false},
	KeyPress{228, false},
	KeyPress{229, false},
	KeyPress{230, false},
	KeyPress{231, false},
	KeyPress{232, false},
	KeyPress{233, false},
	KeyPress{234, false},
	KeyPress{235, false},
	KeyPress{236, false},
	KeyPress{237, false},
	KeyPress{238, false},
	KeyPress{239, false},
	KeyPress{240, false},
	KeyPress{241, false},
	KeyPress{242, false},
	KeyPress{243, false},
	KeyPress{244, false},
	KeyPress{245, false},
	KeyPress{246, false},
	KeyPress{247, false},
	KeyPress{248, false},
	KeyPress{249, false},
	KeyPress{250, false},
	KeyPress{251, false},
	KeyPress{252, false},
	KeyPress{253, false},
	KeyPress{254, false},
	KeyPress{255, false},
	KeyPress{256, false},
	KeyPress{257, false},
	KeyPress{258, false},
	KeyPress{259, false},
	KeyPress{260, false},
	KeyPress{261, false},
	KeyPress{262, false},
	KeyPress{263, false},
	KeyPress{264, false},
	KeyPress{265, false},
	KeyPress{266, false},
	KeyPress{267, false},
	KeyPress{268, false},
	KeyPress{269, false},
	KeyPress{270, false},
	KeyPress{271, false},
	KeyPress{272, false},
	KeyPress{273, false},
	KeyPress{274, false},
	KeyPress{275, false},
	KeyPress{276, false},
	KeyPress{277, false},
	KeyPress{278, false},
	KeyPress{279, false},
	KeyPress{280, false},
	KeyPress{281, false},
	KeyPress{282, false},
	KeyPress{283, false},
	KeyPress{284, false},
	KeyPress{285, false},
	KeyPress{286, false},
	KeyPress{287, false},
	KeyPress{288, false},
	KeyPress{289, false},
	KeyPress{290, false},
	KeyPress{291, false},
	KeyPress{292, false},
	KeyPress{293, false},
	KeyPress{294, false},
	KeyPress{295, false},
	KeyPress{296, false},
	KeyPress{297, false},
	KeyPress{298, false},
	KeyPress{299, false},
	KeyPress{300, false},
	KeyPress{301, false},
	KeyPress{302, false},
	KeyPress{303, false},
	KeyPress{304, false},
	KeyPress{305, false},
	KeyPress{306, false},
	KeyPress{307, false},
	KeyPress{308, false},
	KeyPress{309, false},
	KeyPress{310, false},
	KeyPress{311, false},
	KeyPress{312, false},
	KeyPress{313, false},
	KeyPress{314, false},
	KeyPress{315, false},
	KeyPress{316, false},
	KeyPress{317, false},
	KeyPress{318, false},
	KeyPress{319, false},
	KeyPress{320, false},
	KeyPress{321, false},
	KeyPress{322, false},
	KeyPress{323, false},
	KeyPress{324, false},
	KeyPress{325, false},
	KeyPress{326, false},
	KeyPress{327, false},
	KeyPress{328, false},
	KeyPress{329, false},
	KeyPress{330, false},
	KeyPress{331, false},
	KeyPress{332, false},
	KeyPress{333, false},
	KeyPress{334, false},
	KeyPress{335, false},
	KeyPress{336, false},
	KeyPress{337, false},
	KeyPress{338, false},
	KeyPress{339, false},
	KeyPress{340, false},
	KeyPress{341, false},
	KeyPress{342, false},
	KeyPress{343, false},
	KeyPress{344, false},
	KeyPress{345, false},
	KeyPress{346, false},
	KeyPress{347, false},
	KeyPress{348, false},
	KeyPress{349, false},
	KeyPress{350, false},
	KeyPress{351, false},
	KeyPress{352, false},
	KeyPress{353, false},
	KeyPress{354, false},
	KeyPress{355, false},
	KeyPress{356, false},
	KeyPress{357, false},
	KeyPress{358, false},
	KeyPress{359, false},
	KeyPress{360, false},
	KeyPress{361, false},
	KeyPress{362, false},
	KeyPress{363, false},
	KeyPress{364, false},
	KeyPress{365, false},
	KeyPress{366, false},
	KeyPress{367, false},
	KeyPress{368, false},
	KeyPress{369, false},
	KeyPress{370, false},
	KeyPress{371, false},
	KeyPress{372, false},
	KeyPress{373, false},
	KeyPress{374, false},
	KeyPress{375, false},
	KeyPress{376, false},
	KeyPress{377, false},
	KeyPress{378, false},
	KeyPress{379, false},
	KeyPress{380, false},
	KeyPress{381, false},
	KeyPress{382, false},
	KeyPress{383, false},
	KeyPress{384, false},
	KeyPress{385, false},
	KeyPress{386, false},
	KeyPress{387, false},
	KeyPress{388, false},
	KeyPress{389, false},
	KeyPress{390, false},
	KeyPress{391, false},
	KeyPress{392, false},
	KeyPress{393, false},
	KeyPress{394, false},
	KeyPress{395, false},
	KeyPress{396, false},
	KeyPress{397, false},
	KeyPress{398, false},
	KeyPress{399, false},
	KeyPress{400, false},
	KeyPress{401, false},
	KeyPress{402, false},
	KeyPress{403, false},
	KeyPress{404, false},
	KeyPress{405, false},
	KeyPress{406, false},
	KeyPress{407, false},
	KeyPress{408, false},
	KeyPress{409, false},
	KeyPress{410, false},
	KeyPress{411, false},
	KeyPress{412, false},
	KeyPress{413, false},
	KeyPress{414, false},
	KeyPress{415, false},
	KeyPress{416, false},
	KeyPress{417, false},
	KeyPress{418, false},
	KeyPress{419, false},
	KeyPress{420, false},
	KeyPress{421, false},
	KeyPress{422, false},
	KeyPress{423, false},
	KeyPress{424, false},
	KeyPress{425, false},
	KeyPress{426, false},
	KeyPress{427, false},
	KeyPress{428, false},
	KeyPress{429, false},
	KeyPress{430, false},
	KeyPress{431, false},
	KeyPress{432, false},
	KeyPress{433, false},
	KeyPress{434, false},
	KeyPress{435, false},
	KeyPress{436, false},
	KeyPress{437, false},
	KeyPress{438, false},
	KeyPress{439, false},
	KeyPress{440, false},
	KeyPress{441, false},
	KeyPress{442, false},
	KeyPress{443, false},
	KeyPress{444, false},
	KeyPress{445, false},
	KeyPress{446, false},
	KeyPress{447, false},
	KeyPress{448, false},
	KeyPress{449, false},
	KeyPress{450, false},
	KeyPress{451, false},
	KeyPress{452, false},
	KeyPress{453, false},
	KeyPress{454, false},
	KeyPress{455, false},
	KeyPress{456, false},
	KeyPress{457, false},
	KeyPress{458, false},
	KeyPress{459, false},
	KeyPress{460, false},
	KeyPress{461, false},
	KeyPress{462, false},
	KeyPress{463, false},
	KeyPress{464, false},
	KeyPress{465, false},
	KeyPress{466, false},
	KeyPress{467, false},
	KeyPress{468, false},
	KeyPress{469, false},
	KeyPress{470, false},
	KeyPress{471, false},
	KeyPress{472, false},
	KeyPress{473, false},
	KeyPress{474, false},
	KeyPress{475, false},
	KeyPress{476, false},
	KeyPress{477, false},
	KeyPress{478, false},
	KeyPress{479, false},
	KeyPress{480, false},
	KeyPress{481, false},
	KeyPress{482, false},
	KeyPress{483, false},
	KeyPress{484, false},
	KeyPress{485, false},
	KeyPress{486, false},
	KeyPress{487, false},
	KeyPress{488, false},
	KeyPress{489, false},
	KeyPress{490, false},
	KeyPress{491, false},
	KeyPress{492, false},
	KeyPress{493, false},
	KeyPress{494, false},
	KeyPress{495, false},
	KeyPress{496, false},
	KeyPress{497, false},
	KeyPress{498, false},
	KeyPress{499, false},
	KeyPress{500, false},
	KeyPress{501, false},
	KeyPress{502, false},
	KeyPress{503, false},
	KeyPress{504, false},
	KeyPress{505, false},
	KeyPress{506, false},
	KeyPress{507, false},
	KeyPress{508, false},
	KeyPress{509, false},
	KeyPress{510, false},
	KeyPress{511, false},
	KeyPress{512, false},
	KeyPress{513, false},
	KeyPress{514, false},
	KeyPress{515, false},
	KeyPress{516, false},
	KeyPress{517, false},
	KeyPress{518, false},
	KeyPress{519, false},
	KeyPress{520, false},
	KeyPress{521, false},
	KeyPress{522, false},
	KeyPress{523, false},
	KeyPress{524, false},
	KeyPress{525, false},
	KeyPress{526, false},
	KeyPress{527, false},
	KeyPress{528, false},
	KeyPress{529, false},
	KeyPress{530, false},
	KeyPress{531, false},
	KeyPress{532, false},
	KeyPress{533, false},
	KeyPress{534, false},
	KeyPress{535, false},
	KeyPress{536, false},
	KeyPress{537, false},
	KeyPress{538, false},
	KeyPress{539, false},
	KeyPress{540, false},
	KeyPress{541, false},
	KeyPress{542, false},
	KeyPress{543, false},
	KeyPress{544, false},
	KeyPress{545, false},
	KeyPress{546, false},
	KeyPress{547, false},
	KeyPress{548, false},
	KeyPress{549, false},
	KeyPress{550, false},
	KeyPress{551, false},
	KeyPress{552, false},
	KeyPress{553, false},
	KeyPress{554, false},
	KeyPress{555, false},
	KeyPress{556, false},
	KeyPress{557, false},
	KeyPress{558, false},
	KeyPress{559, false},
	KeyPress{560, false},
	KeyPress{561, false},
	KeyPress{562, false},
	KeyPress{563, false},
	KeyPress{564, false},
	KeyPress{565, false},
	KeyPress{566, false},
	KeyPress{567, false},
	KeyPress{568, false},
	KeyPress{569, false},
	KeyPress{570, false},
	KeyPress{571, false},
	KeyPress{572, false},
	KeyPress{573, false},
	KeyPress{574, false},
	KeyPress{575, false},
	KeyPress{576, false},
	KeyPress{577, false},
	KeyPress{578, false},
	KeyPress{579, false},
	KeyPress{580, false},
	KeyPress{581, false},
	KeyPress{582, false},
	KeyPress{583, false},
	KeyPress{584, false},
	KeyPress{585, false},
	KeyPress{586, false},
	KeyPress{587, false},
	KeyPress{588, false},
	KeyPress{589, false},
	KeyPress{590, false},
	KeyPress{591, false},
	KeyPress{592, false},
	KeyPress{593, false},
	KeyPress{594, false},
	KeyPress{595, false},
	KeyPress{596, false},
	KeyPress{597, false},
	KeyPress{598, false},
	KeyPress{599, false},
	KeyPress{600, false},
	KeyPress{601, false},
	KeyPress{602, false},
	KeyPress{603, false},
	KeyPress{604, false},
	KeyPress{605, false},
	KeyPress{606, false},
	KeyPress{607, false},
	KeyPress{608, false},
	KeyPress{609, false},
	KeyPress{610, false},
	KeyPress{611, false},
	KeyPress{612, false},
	KeyPress{613, false},
	KeyPress{614, false},
	KeyPress{615, false},
	KeyPress{616, false},
	KeyPress{617, false},
	KeyPress{618, false},
	KeyPress{619, false},
	KeyPress{620, false},
	KeyPress{621, false},
	KeyPress{622, false},
	KeyPress{623, false},
	KeyPress{624, false},
	KeyPress{625, false},
	KeyPress{626, false},
	KeyPress{627, false},
	KeyPress{628, false},
	KeyPress{629, false},
	KeyPress{630, false},
	KeyPress{631, false},
	KeyPress{632, false},
	KeyPress{633, false},
	KeyPress{634, false},
	KeyPress{635, false},
	KeyPress{636, false},
	KeyPress{637, false},
	KeyPress{638, false},
	KeyPress{639, false},
	KeyPress{640, false},
	KeyPress{641, false},
	KeyPress{642, false},
	KeyPress{643, false},
	KeyPress{644, false},
	KeyPress{645, false},
	KeyPress{646, false},
	KeyPress{647, false},
	KeyPress{648, false},
	KeyPress{649, false},
	KeyPress{650, false},
	KeyPress{651, false},
	KeyPress{652, false},
	KeyPress{653, false},
	KeyPress{654, false},
	KeyPress{655, false},
	KeyPress{656, false},
	KeyPress{657, false},
	KeyPress{658, false},
	KeyPress{659, false},
	KeyPress{660, false},
	KeyPress{661, false},
	KeyPress{662, false},
	KeyPress{663, false},
	KeyPress{664, false},
	KeyPress{665, false},
	KeyPress{666, false},
	KeyPress{667, false},
	KeyPress{668, false},
	KeyPress{669, false},
	KeyPress{670, false},
	KeyPress{671, false},
	KeyPress{672, false},
	KeyPress{673, false},
	KeyPress{674, false},
	KeyPress{675, false},
	KeyPress{676, false},
	KeyPress{677, false},
	KeyPress{678, false},
	KeyPress{679, false},
	KeyPress{680, false},
	KeyPress{681, false},
	KeyPress{682, false},
	KeyPress{683, false},
	KeyPress{684, false},
	KeyPress{685, false},
	KeyPress{686, false},
	KeyPress{687, false},
	KeyPress{688, false},
	KeyPress{689, false},
	KeyPress{690, false},
	KeyPress{691, false},
	KeyPress{692, false},
	KeyPress{693, false},
	KeyPress{694, false},
	KeyPress{695, false},
	KeyPress{696, false},
	KeyPress{697, false},
	KeyPress{698, false},
	KeyPress{699, false},
	KeyPress{700, false},
	KeyPress{701, false},
	KeyPress{702, false},
	KeyPress{703, false},
	KeyPress{704, false},
	KeyPress{705, false},
	KeyPress{706, false},
	KeyPress{707, false},
	KeyPress{708, false},
	KeyPress{709, false},
	KeyPress{710, false},
	KeyPress{711, false},
	KeyPress{712, false},
	KeyPress{713, false},
	KeyPress{714, false},
	KeyPress{715, false},
	KeyPress{716, false},
	KeyPress{717, false},
	KeyPress{718, false},
	KeyPress{719, false},
	KeyPress{720, false},
	KeyPress{721, false},
	KeyPress{722, false},
	KeyPress{723, false},
	KeyPress{724, false},
	KeyPress{725, false},
	KeyPress{726, false},
	KeyPress{727, false},
	KeyPress{728, false},
	KeyPress{729, false},
	KeyPress{730, false},
	KeyPress{731, false},
	KeyPress{732, false},
	KeyPress{733, false},
	KeyPress{734, false},
	KeyPress{735, false},
	KeyPress{736, false},
	KeyPress{737, false},
	KeyPress{738, false},
	KeyPress{739, false},
	KeyPress{740, false},
	KeyPress{741, false},
	KeyPress{742, false},
	KeyPress{743, false},
	KeyPress{744, false},
	KeyPress{745, false},
	KeyPress{746, false},
	KeyPress{747, false},
	KeyPress{748, false},
	KeyPress{749, false},
	KeyPress{750, false},
	KeyPress{751, false},
	KeyPress{752, false},
	KeyPress{753, false},
	KeyPress{754, false},
	KeyPress{755, false},
	KeyPress{756, false},
	KeyPress{757, false},
	KeyPress{758, false},
	KeyPress{759, false},
	KeyPress{760, false},
	KeyPress{761, false},
	KeyPress{762, false},
	KeyPress{763, false},
	KeyPress{764, false},
	KeyPress{765, false},
	KeyPress{766, false},
	KeyPress{767, false},
	KeyPress{768, false},
	KeyPress{769, false},
	KeyPress{770, false},
	KeyPress{771, false},
	KeyPress{772, false},
	KeyPress{773, false},
	KeyPress{774, false},
	KeyPress{775, false},
	KeyPress{776, false},
	KeyPress{777, false},
	KeyPress{778, false},
	KeyPress{779, false},
	KeyPress{780, false},
	KeyPress{781, false},
	KeyPress{782, false},
	KeyPress{783, false},
	KeyPress{784, false},
	KeyPress{785, false},
	KeyPress{786, false},
	KeyPress{787, false},
	KeyPress{788, false},
	KeyPress{789, false},
	KeyPress{790, false},
	KeyPress{791, false},
	KeyPress{792, false},
	KeyPress{793, false},
	KeyPress{794, false},
	KeyPress{795, false},
	KeyPress{796, false},
	KeyPress{797, false},
	KeyPress{798, false},
	KeyPress{799, false},
	KeyPress{800, false},
	KeyPress{801, false},
	KeyPress{802, false},
	KeyPress{803, false},
	KeyPress{804, false},
	KeyPress{805, false},
	KeyPress{806, false},
	KeyPress{807, false},
	KeyPress{808, false},
	KeyPress{809, false},
	KeyPress{810, false},
	KeyPress{811, false},
	KeyPress{812, false},
	KeyPress{813, false},
	KeyPress{814, false},
	KeyPress{815, false},
	KeyPress{816, false},
	KeyPress{817, false},
	KeyPress{818, false},
	KeyPress{819, false},
	KeyPress{820, false},
	KeyPress{821, false},
	KeyPress{822, false},
	KeyPress{823, false},
	KeyPress{824, false},
	KeyPress{825, false},
	KeyPress{826, false},
	KeyPress{827, false},
	KeyPress{828, false},
	KeyPress{829, false},
	KeyPress{830, false},
	KeyPress{831, false},
	KeyPress{832, false},
	KeyPress{833, false},
	KeyPress{834, false},
	KeyPress{835, false},
	KeyPress{836, false},
	KeyPress{837, false},
	KeyPress{838, false},
	KeyPress{839, false},
	KeyPress{840, false},
	KeyPress{841, false},
	KeyPress{842, false},
	KeyPress{843, false},
	KeyPress{844, false},
	KeyPress{845, false},
	KeyPress{846, false},
	KeyPress{847, false},
	KeyPress{848, false},
	KeyPress{849, false},
	KeyPress{850, false},
	KeyPress{851, false},
	KeyPress{852, false},
	KeyPress{853, false},
	KeyPress{854, false},
	KeyPress{855, false},
	KeyPress{856, false},
	KeyPress{857, false},
	KeyPress{858, false},
	KeyPress{859, false},
	KeyPress{860, false},
	KeyPress{861, false},
	KeyPress{862, false},
	KeyPress{863, false},
	KeyPress{864, false},
	KeyPress{865, false},
	KeyPress{866, false},
	KeyPress{867, false},
	KeyPress{868, false},
	KeyPress{869, false},
	KeyPress{870, false},
	KeyPress{871, false},
	KeyPress{872, false},
	KeyPress{873, false},
	KeyPress{874, false},
	KeyPress{875, false},
	KeyPress{876, false},
	KeyPress{877, false},
	KeyPress{878, false},
	KeyPress{879, false},
	KeyPress{880, false},
	KeyPress{881, false},
	KeyPress{882, false},
	KeyPress{883, false},
	KeyPress{884, false},
	KeyPress{885, false},
	KeyPress{886, false},
	KeyPress{887, false},
	KeyPress{888, false},
	KeyPress{889, false},
	KeyPress{890, false},
	KeyPress{891, false},
	KeyPress{892, false},
	KeyPress{893, false},
	KeyPress{894, false},
	KeyPress{895, false},
	KeyPress{896, false},
	KeyPress{897, false},
	KeyPress{898, false},
	KeyPress{899, false},
	KeyPress{900, false},
	KeyPress{901, false},
	KeyPress{902, false},
	KeyPress{903, false},
	KeyPress{904, false},
	KeyPress{905, false},
	KeyPress{906, false},
	KeyPress{907, false},
	KeyPress{908, false},
	KeyPress{909, false},
	KeyPress{910, false},
	KeyPress{911, false},
	KeyPress{912, false},
	KeyPress{913, false},
	KeyPress{914, false},
	KeyPress{915, false},
	KeyPress{916, false},
	KeyPress{917, false},
	KeyPress{918, false},
	KeyPress{919, false},
	KeyPress{920, false},
	KeyPress{921, false},
	KeyPress{922, false},
	KeyPress{923, false},
	KeyPress{924, false},
	KeyPress{925, false},
	KeyPress{926, false},
	KeyPress{927, false},
	KeyPress{928, false},
	KeyPress{929, false},
	KeyPress{930, false},
	KeyPress{931, false},
	KeyPress{932, false},
	KeyPress{933, false},
	KeyPress{934, false},
	KeyPress{935, false},
	KeyPress{936, false},
	KeyPress{937, false},
	KeyPress{938, false},
	KeyPress{939, false},
	KeyPress{940, false},
	KeyPress{941, false},
	KeyPress{942, false},
	KeyPress{943, false},
	KeyPress{944, false},
	KeyPress{945, false},
	KeyPress{946, false},
	KeyPress{947, false},
	KeyPress{948, false},
	KeyPress{949, false},
	KeyPress{950, false},
	KeyPress{951, false},
	KeyPress{952, false},
	KeyPress{953, false},
	KeyPress{954, false},
	KeyPress{955, false},
	KeyPress{956, false},
	KeyPress{957, false},
	KeyPress{958, false},
	KeyPress{959, false},
	KeyPress{960, false},
	KeyPress{961, false},
	KeyPress{962, false},
	KeyPress{963, false},
	KeyPress{964, false},
	KeyPress{965, false},
	KeyPress{966, false},
	KeyPress{967, false},
	KeyPress{968, false},
	KeyPress{969, false},
	KeyPress{970, false},
	KeyPress{971, false},
	KeyPress{972, false},
	KeyPress{973, false},
	KeyPress{974, false},
	KeyPress{975, false},
	KeyPress{976, false},
	KeyPress{977, false},
	KeyPress{978, false},
	KeyPress{979, false},
	KeyPress{980, false},
	KeyPress{981, false},
	KeyPress{982, false},
	KeyPress{983, false},
	KeyPress{984, false},
	KeyPress{985, false},
	KeyPress{986, false},
	KeyPress{987, false},
	KeyPress{988, false},
	KeyPress{989, false},
	KeyPress{990, false},
	KeyPress{991, false},
	KeyPress{992, false},
	KeyPress{993, false},
	KeyPress{994, false},
	KeyPress{995, false},
	KeyPress{996, false},
	KeyPress{997, false},
	KeyPress{998, false},
	KeyPress{999, false},
	KeyPress{1000, false},
	KeyPress{1001, false},
	KeyPress{1002, false},
	KeyPress{1003, false},
	KeyPress{1004, false},
	KeyPress{1005, false},
	KeyPress{1006, false},
	KeyPress{1007, false},
	KeyPress{1008, false},
	KeyPress{1009, false},
	KeyPress{1010, false},
	KeyPress{1011, false},
	KeyPress{1012, false},
	KeyPress{1013, false},
	KeyPress{1014, false},
	KeyPress{1015, false},
	KeyPress{1016, false},
	KeyPress{1017, false},
	KeyPress{1018, false},
	KeyPress{1019, false},
	KeyPress{1020, false},
	KeyPress{1021, false},
	KeyPress{1022, false},
	KeyPress{1023, false},
}
