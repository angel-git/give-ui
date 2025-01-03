package cache

import "testing"

func TestGetMaxVisibleAmmoRanges(t *testing.T) {
	if len(getMaxVisibleAmmoRanges("1-1;4-4;6-6;8-8;10-10;12-12;14-14;16-16;18-18;20-20")) != 10 {
		t.Error("getMaxVisibleAmmoRanges() failed")
	}
	if len(getMaxVisibleAmmoRanges("1-3")) != 1 {
		t.Error("getMaxVisibleAmmoRanges() failed")
	}
}

func TestGetMaxVisibleAmmo(t *testing.T) {
	if getMaxVisibleAmmo(0, "1-1;4-4;6-6;8-8;10-10;12-12;14-14;16-16;18-18;20-20") != 1 {
		t.Error("getMaxVisibleAmmo() failed")
	}
	if getMaxVisibleAmmo(1, "1-1;4-4;6-6;8-8;10-10;12-12;14-14;16-16;18-18;20-20") != 1 {
		t.Error("getMaxVisibleAmmo() failed")
	}
	if getMaxVisibleAmmo(2, "1-1;4-4;6-6;8-8;10-10;12-12;14-14;16-16;18-18;20-20") != 1 {
		t.Error("getMaxVisibleAmmo() failed")
	}
	if getMaxVisibleAmmo(3, "1-1;4-4;6-6;8-8;10-10;12-12;14-14;16-16;18-18;20-20") != 1 {
		t.Error("getMaxVisibleAmmo() failed")
	}
	if getMaxVisibleAmmo(4, "1-1;4-4;6-6;8-8;10-10;12-12;14-14;16-16;18-18;20-20") != 4 {
		t.Error("getMaxVisibleAmmo() failed")
	}
	if getMaxVisibleAmmo(20, "1-1;4-4;6-6;8-8;10-10;12-12;14-14;16-16;18-18;20-20") != 20 {
		t.Error("getMaxVisibleAmmo() failed")
	}
}

func TestGetDeterministicHashCode(t *testing.T) {
	if getDeterministicHashCode("5448ba0b4bdc2d02308b456c") != 1091773418 {
		t.Error("get_deterministic_hash_code() failed")
	}
}

func TestGetHashCodeFromMongoID(t *testing.T) {
	if getHashCodeFromMongoID("5448ba0b4bdc2d02308b456c") != -865571915 {
		t.Error("get_hash_code_from_mongo_id() failed")
	}
}
