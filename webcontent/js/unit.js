/** 
 * @constructor
 * @param{crisisJson.Unit} unitJson 
 * @param{crisis.DivisionDetails} dets
 */ 
crisis.Unit = function(unitJson, dets) {
    var unit = this;
    
    /** @type{crisis.DivisionDetails} */
    unit.details = dets;
    /** @type{number} */
    unit.amount = unitJson.Amount;
    /** @type{string} */
    unit.typeName = unitJson.TypeName;
    /** @type{number} */
    unit.typeNum = unitJson.TypeNum;
    /** @type{jQuery} */
    unit.$listItem = crisis.cloneProto(crisis.$protoUnitListItem);
    /** @type{jQuery} */
    unit.$value = unit.$listItem.find(".value").html(unit.amount);
    /** @type{jQuery} */
    unit.$type = unit.$listItem.find(".type").html(crisis.Unit.typeHtml(unit.typeNum));
    /** @type{jQuery} */
    unit.$editField = unit.$listItem.find(".editField");
    /** @type{jQuery} */
    unit.$removeUnitButton = unit.$listItem.find(".removeUnitButton");
    unit.$removeUnitButton.on("click.crisis", function() {
        unit.details.removeUnit(unit);
    });
    /** @type{jQuery} */
    unit.$invalidAlert = unit.$listItem.find(".invalidAlert");
}

crisis.Unit.prototype.editOn = function() {
    this.$editField.val(this.amount).show();
    this.$removeUnitButton.show();
    this.$value.hide();
}

crisis.Unit.prototype.editOff = function() {
    this.$editField.hide();
    this.$removeUnitButton.hide();
    this.$value.show();
}

/** @param{number} typeId */
crisis.Unit.typeHtml = function(typeId) {
    return crisis.cloneProto(
	      crisis.$protoUnitTypes.find(crisis.unitTypeSelector(typeId)));
}

/** @param{crisisJson.Unit} unitJson */
crisis.Unit.fromData = function(unitJson) {
    return new crisis.Unit(unitJson);
}
