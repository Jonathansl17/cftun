package routes

import "context"

func (c DNSCleaner) Delete(ctx context.Context, host string) (DNSRemoval, error) {
	configured, err := c.API.Configured()
	if err != nil {
		return DNSRemoval{}, err
	}
	if !configured {
		return DNSRemoval{State: DNSManual}, nil
	}
	count, err := c.API.DeleteCNAME(ctx, host)
	removal := DNSRemoval{State: DNSDeleted, Deleted: count}
	if err != nil && count == 0 {
		removal.State = DNSNotAttempted
	}
	return removal, err
}
