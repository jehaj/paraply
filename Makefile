D2_FILES := $(wildcard docs/*.d2)
SVG_FILES := $(D2_FILES:.d2=.svg)

docs: $(SVG_FILES)

docs/%.svg: docs/%.d2
	d2 $< $@

.PHONY: docs
