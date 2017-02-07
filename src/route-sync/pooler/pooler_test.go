package pooler_test

import (
	"route-sync/pooler"
	"route-sync/route"
	routefakes "route-sync/route/fakes"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Time-based Pooler", func() {

	It("starts and stops", func() {
		pool := pooler.ByTime(time.Duration(0))

		running := func() bool { return pool.Running() }

		done, _ := pool.Start(&routefakes.Source{}, &routefakes.Sink{})
		Eventually(running).Should(BeTrue())
		done <- struct{}{}
		Eventually(running).Should(BeFalse())
	})

	It("pools and passes", func() {
		pool := pooler.ByTime(time.Duration(0))

		src := &routefakes.Source{}
		tcpRoute := &route.TCP{Frontend: 8080,
			Backend: []route.Endpoint{route.Endpoint{IP: "10.10.0.10", Port: 9090}}}
		httpRoute := &route.HTTP{Name: "foo.bar.com",
			Backend: []route.Endpoint{route.Endpoint{IP: "10.10.0.10", Port: 9090}}}

		src.TCP_value = []*route.TCP{tcpRoute}
		src.HTTP_value = []*route.HTTP{httpRoute}
		sink := &routefakes.Sink{}

		done, _ := pool.Start(src, sink)

		Eventually(func() bool { return src.TCP_count > 0 }).Should(BeTrue())
		Eventually(func() bool { return sink.TCP_count > 0 }).Should(BeTrue())
		Expect(sink.TCP_values[0][0]).To(Equal(tcpRoute))

		Eventually(func() bool { return src.HTTP_count > 0 }).Should(BeTrue())
		Eventually(func() bool { return sink.HTTP_count > 0 }).Should(BeTrue())
		Expect(sink.HTTP_values[0][0]).To(Equal(httpRoute))

		done <- struct{}{}
	})
})