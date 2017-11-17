package main

var schema = `
DROP TABLE IF EXISTS cell;
CREATE TABLE cell (
    ts               DATETIME(6),
    SetVoltage       DECIMAL(6,3),
	SetCurrent       DECIMAL(6,3),
	ActualVoltage    DECIMAL(6,3),
	ActualCurrent    DECIMAL(6,3),
	Power            DECIMAL(6,3),
	SupplyVoltage    DECIMAL(6,3),
	ProtectionTrip   TINYINT,
	Constant         TINYINT,
	OutputOn         TINYINT,
    PRIMARY KEY (ts)
);`
