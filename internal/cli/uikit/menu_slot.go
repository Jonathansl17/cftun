package uikit

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func MarkMenu(cmd *cobra.Command, slot MenuSlot) *cobra.Command {
	setAnnotation(cmd, annotationMenuGroup, string(slot.Group))
	setAnnotation(cmd, annotationMenuLabel, slot.Label)
	setAnnotation(cmd, annotationMenuOrder, strconv.Itoa(slot.Order))
	return cmd
}

func MarkIn(group MenuGroup, cmd *cobra.Command, label string, order int) *cobra.Command {
	return MarkMenu(cmd, MenuSlot{Group: group, Label: label, Order: order})
}

func SlotOf(cmd *cobra.Command) (MenuSlot, bool, error) {
	label, marked := cmd.Annotations[annotationMenuLabel]
	if !marked {
		return MenuSlot{}, false, nil
	}
	order, err := strconv.Atoi(cmd.Annotations[annotationMenuOrder])
	if err != nil {
		return MenuSlot{}, false, fmt.Errorf(slotErrorFormat, cmd.Name(), err)
	}
	group := MenuGroup(cmd.Annotations[annotationMenuGroup])
	return MenuSlot{Group: group, Label: label, Order: order}, true, nil
}
