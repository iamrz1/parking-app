package model

import cerror "supertal-tha-parking-app/error"

func (u *UserCreateReq) Validate() error {
	err := cerror.ValidationError{}
	if u.Name == "" {
		err.Add("message", "Name is required")
	}

	if u.Username == "" {
		err.Add("message", "Username is required")
	}

	if u.Password == "" {
		err.Add("message", "Password is required")
	}

	if len(err) != 0 {
		return err
	}

	return nil
}

func (u *LoginReq) Validate() error {
	err := cerror.ValidationError{}

	if u.Username == "" {
		err.Add("message", "Username is required")
	}

	if u.Password == "" {
		err.Add("message", "Password is required")
	}

	if len(err) != 0 {
		return err
	}

	return nil
}

func (b *BookingReq) Validate() error {
	err := cerror.ValidationError{}

	if b.LotID == 0 {
		err.Add("message", "Lot ID is required")
	}

	if len(err) != 0 {
		return err
	}

	return nil
}

func (p *ParkingLotCreateReq) Validate() error {
	err := cerror.ValidationError{}

	if p.Name == "" {
		err.Add("message", "Name is required")
	}

	if p.NumberOfSlots == 0 {
		err.Add("message", "NumberOfSlots is required")
	}

	if len(err) != 0 {
		return err
	}

	return nil
}

func (m *MaintenanceStatusReq) Validate() error {
	err := cerror.ValidationError{}

	if m.MaintenanceMode == nil {
		err.Add("message", "MaintenanceMode is required")
	}

	if m.SlotID == 0 {
		err.Add("message", "SlotID ID is required")
	}

	if len(err) != 0 {
		return err
	}

	return nil
}
