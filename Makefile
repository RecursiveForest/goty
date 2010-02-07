include $(GOROOT)/src/Make.$(GOARCH)
 
TARG=goty
GOFILES=goty.go

CLEANFILES+=sic

include $(GOROOT)/src/Make.pkg 

sic: install sic.go
	$(QUOTED_GOBIN)/$(GC) sic.go
	$(QUOTED_GOBIN)/$(LD) -o sic sic.6
