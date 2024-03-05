package models

type Zsd001nModels struct {
	ID               int    `json:"id"`
	BillingDocument  string `json:"billingDocument"`
	BillingType      string `json:"billingType"`
	DistrictName     string `json:"districtName"`
	Order            string `json:"order"`
	OrderReasonDesc  string `json:"orderReasonDesc"`
	BillingDate      string `json:"billingDate"`
	CreatedOn        string `json:"createdOn"`
	Slor             string `json:"slor"`
	Payer            string `json:"payer"`
	PayerName        string `json:"payerName"`
	Dc               string `json:"Dc"`
	Promotion        string `json:"promotion"`
	PromoDesc        string `json:"promoDesc"`
	SProm            string `json:"sProm"`
	ProId            string `json:"proId"`
	LocCode          string `json:"locCode"`
	MaterialNumber   string `json:"materialNumber"`
	MaterialDesc     string `json:"materialDesc"`
	Item             string `json:"item"`
	BillingQty       string `json:"billingQty"`
	DppLc            string `json:"dppLc"`
	PpnLc            string `json:"ppnLc"`
	Total            string `json:"total"`
	SlsDoc           string `json:"slsDoc"`
	AccDoc           string `json:"accDoc"`
	Description      string `json:"description"`
	Mg1              string `json:"mg1"`
	MatGrpDesc1      string `json:"matGrpDesc1"`
	Mg2              string `json:"mg2"`
	MatGrpDesc2      string `json:"matGrpDesc2"`
	Mg3              string `json:"mg3"`
	MatGrpDesc3      string `json:"matGrpDesc3"`
	Mg4              string `json:"mg4"`
	MatGrpDesc4      string `json:"matGrpDesc4"`
	Mg5              string `json:"mg5"`
	MatGrpDesc5      string `json:"matGrpDesc5"`
	Mg6              string `json:"mg6"`
	MatGrpDesc6      string `json:"matGrpDesc6"`
	Mg7              string `json:"mg7"`
	MatGrpDesc7      string `json:"matGrpDesc7"`
	Mg8              string `json:"mg8"`
	MatGrpDesc8      string `json:"matGrpDesc8"`
	Mg9              string `json:"mg9"`
	MatGrpDesc9      string `json:"matGrpDesc9"`
	Cg               string `json:"cg"`
	CustGrpDesc      string `json:"custGrpDesc"`
	IndCode1         string `json:"indCode1"`
	Indry            string `json:"indry"`
	DescriptionIndry string `json:"descriptionIndry"`
	MatGrp           string `json:"matGrp"`
	Lab              string `json:"lab"`
	LabDesc          string `json:"labDesc"`
	ShipPointDesc    string `json:"shipPointDesc"`
	ShPo             string `json:"shPo"`
	Created          string `json:"created"`
	PoNumber         string `json:"poNumber"`
	Sgr              string `json:"sgr"`
	SlsGrpDesc       string `json:"slsGrpDesc"`
	ShCss            string `json:"shCss"`
	ShCssDesc        string `json:"shCssDesc"`
	NoFakturPajak    string `json:"noFakturPajak"`
	ShipToAddress    string `json:"shipToAddress"`
	StatusCancelled  string `json:"statusCancelled"`
	DcDesc           string `json:"dcDesc"`
	ShipToPa         string `json:"shipToPa"`
	ShipToPartyDesc  string `json:"shipToPartyDesc"`
	Reference        string `json:"reference"`
	Plnt             string `json:"plnt"`
	Sloc             string `json:"sloc"`
	Aag              string `json:"aag"`
	AssNumber        string `json:"assNumber"`
	BillCancelledRef string `json:"billCancelledRef"`
}

func (Zsd001nModels) TableName() string {
	return "sf_dump_zsd001n"
}
