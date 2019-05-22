package statsd

import (
	"strconv"
	"strings"
)

var (
	gaugeSymbol        = []byte("g")
	countSymbol        = []byte("c")
	histogramSymbol    = []byte("h")
	distributionSymbol = []byte("d")
	setSymbol          = []byte("s")
	timingSymbol       = []byte("ms")
)

func appendHeader(buffer []byte, namespace string, name string) []byte {
	if namespace != "" {
		buffer = append(buffer, namespace...)
	}
	buffer = append(buffer, name...)
	buffer = append(buffer, ':')
	return buffer
}

func appendRate(buffer []byte, rate float64) []byte {
	if rate < 1 {
		buffer = append(buffer, "|@"...)
		buffer = strconv.AppendFloat(buffer, rate, 'f', -1, 64)
	}
	return buffer
}

func appendWithoutNewlines(buffer []byte, s string) []byte {
	// fastpath for strings without newlines
	if strings.IndexByte(s, '\n') == -1 {
		return append(buffer, s...)
	}

	for _, b := range []byte(s) {
		if b != '\n' {
			buffer = append(buffer, b)
		}
	}
	return buffer
}

func appendTags(buffer []byte, globalTags []string, tags []string) []byte {
	buffer = append(buffer, "|#"...)
	firstTag := true

	for _, tag := range globalTags {
		if !firstTag {
			buffer = append(buffer, ',')
		}
		buffer = appendWithoutNewlines(buffer, tag)
		firstTag = false
	}
	for _, tag := range tags {
		if !firstTag {
			buffer = append(buffer, ',')
		}
		buffer = appendWithoutNewlines(buffer, tag)
		firstTag = false
	}
	return buffer
}

func appendFloatMetric(buffer []byte, typeSymbol []byte, namespace string, globalTags []string, name string, value float64, tags []string, rate float64) []byte {
	buffer = appendHeader(buffer, namespace, name)
	buffer = strconv.AppendFloat(buffer, value, 'f', 6, 64)
	buffer = append(buffer, '|')
	buffer = append(buffer, typeSymbol...)
	buffer = appendTags(buffer, globalTags, tags)
	return buffer
}

func appendIntegerMetric(buffer []byte, typeSymbol []byte, namespace string, globalTags []string, name string, value int64, tags []string, rate float64) []byte {
	buffer = appendHeader(buffer, namespace, name)
	buffer = strconv.AppendInt(buffer, value, 10)
	buffer = append(buffer, '|')
	buffer = append(buffer, typeSymbol...)
	buffer = appendTags(buffer, globalTags, tags)
	return buffer
}

func appendStringMetric(buffer []byte, typeSymbol []byte, namespace string, globalTags []string, name string, value string, tags []string, rate float64) []byte {
	buffer = appendHeader(buffer, namespace, name)
	buffer = append(buffer, value...)
	buffer = append(buffer, '|')
	buffer = append(buffer, typeSymbol...)
	buffer = appendTags(buffer, globalTags, tags)
	return buffer
}

func appendGauge(buffer []byte, namespace string, globalTags []string, name string, value float64, tags []string, rate float64) []byte {
	c := Client{
		Tags:      globalTags,
		Namespace: namespace,
	}
	suffix := gaugeSuffix
	return c.format(name, value, suffix, tags, rate)
}

func appendCount(buffer []byte, namespace string, globalTags []string, name string, value int64, tags []string, rate float64) []byte {
	c := Client{
		Tags:      globalTags,
		Namespace: namespace,
	}
	suffix := countSuffix
	return c.format(name, value, suffix, tags, rate)
}

func appendHistogram(buffer []byte, namespace string, globalTags []string, name string, value float64, tags []string, rate float64) []byte {
	c := Client{
		Tags:      globalTags,
		Namespace: namespace,
	}
	suffix := histogramSuffix
	return c.format(name, value, suffix, tags, rate)
}

func appendDistribution(buffer []byte, namespace string, globalTags []string, name string, value float64, tags []string, rate float64) []byte {
	c := Client{
		Tags:      globalTags,
		Namespace: namespace,
	}
	suffix := distributionSuffix
	return c.format(name, value, suffix, tags, rate)
}

func appendDecrement(buffer []byte, namespace string, globalTags []string, name string, tags []string, rate float64) []byte {
	c := Client{
		Tags:      globalTags,
		Namespace: namespace,
	}
	suffix := decrSuffix
	return c.format(name, nil, suffix, tags, rate)
}

func appendIncrement(buffer []byte, namespace string, globalTags []string, name string, tags []string, rate float64) []byte {
	c := Client{
		Tags:      globalTags,
		Namespace: namespace,
	}
	suffix := incrSuffix
	return c.format(name, nil, suffix, tags, rate)
}

func appendSet(buffer []byte, namespace string, globalTags []string, name string, value string, tags []string, rate float64) []byte {
	c := Client{
		Tags:      globalTags,
		Namespace: namespace,
	}
	suffix := setSuffix
	return c.format(name, value, suffix, tags, rate)
}

func appendTiming(buffer []byte, namespace string, globalTags []string, name string, value float64, tags []string, rate float64) []byte {
	c := Client{
		Tags:      globalTags,
		Namespace: namespace,
	}
	suffix := timingSuffix
	return c.format(name, value, suffix, tags, rate)
}
