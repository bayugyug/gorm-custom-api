package services_test

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/bayugyug/gorm-custom-api/api/routes"
	"github.com/bayugyug/gorm-custom-api/services"
	"github.com/bayugyug/gorm-custom-api/tools"

	"github.com/icrowley/fake"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var svcrouter *routes.APIRouter
var store *gorm.DB
var service *services.BuildingService

var _ = BeforeSuite(func() {
	var err error
	//init
	svcrouter, err = routes.NewAPIRouter(
		routes.WithSvcOptConfig(tools.DevelTester{}.Config()),
	)
	Expect(err).NotTo(HaveOccurred())
	tools.DevelTester{}.EmptyDBTst(svcrouter.DBHandle.GetConnection())
	//globally
	store = svcrouter.DBHandle.GetConnection()
	service = services.NewBuildingService()
})

var _ = AfterSuite(func() {
	tools.DevelTester{}.EmptyDBTst(svcrouter.DBHandle.GetConnection())
})

var _ = Describe("REST Building API Service::MODELS", func() {

	BeforeEach(func() {
	})

	Context("Valid parameters", func() {

		Context("Create record", func() {
			It("should return ok", func() {
				name := fmt.Sprintf("building::%s", fake.DigitsN(15))
				pid, err := service.Create(store,
					&services.BuildingCreateParams{
						Name:    &name,
						Address: fmt.Sprintf("Marina Boulevard::%s", fake.DigitsN(15)),
						Floors:  tools.Seeder{}.CreateFloors(),
					},
				)
				if err != nil {
					Fail(err.Error())
				}
				Expect(pid).Should(BeNumerically(">", 0))
				By("Create data ok")

			})
		})

		Context("Update record", func() {
			It("should return ok", func() {
				name := fmt.Sprintf("marina-bay-sands::%s", fake.DigitsN(15))
				pid, err := service.Create(store,
					&services.BuildingCreateParams{
						Name:    &name,
						Address: fmt.Sprintf("Marina Boulevard::%s", fake.DigitsN(15)),
						Floors:  tools.Seeder{}.CreateFloors(),
					},
				)
				if err != nil {
					Fail(err.Error())
				}
				Expect(pid).Should(BeNumerically(">", 0))
				By("Create data before update ok")

				err = service.Update(store,
					&services.BuildingUpdateParams{
						ID: &pid,
						BuildingCreateParams: services.BuildingCreateParams{
							Name:    &name,
							Address: fmt.Sprintf("Marina Boulevard::%s", fake.DigitsN(15)),
							Floors:  tools.Seeder{}.CreateFloors(),
						},
					},
				)
				if err != nil {
					Fail(err.Error())
				}
				By("Update data ok")
			})
		})

		Context("Delete record", func() {
			It("should return ok", func() {
				name := fmt.Sprintf("marina-bay-sands::%s", fake.DigitsN(15))
				pid, err := service.Create(store,
					&services.BuildingCreateParams{
						Name:    &name,
						Address: fmt.Sprintf("Marina Boulevard::%s", fake.DigitsN(15)),
						Floors:  tools.Seeder{}.CreateFloors(),
					},
				)
				if err != nil {
					Fail(err.Error())
				}
				Expect(pid).Should(BeNumerically(">", 0))
				By("Create data before delete ok")

				err = service.Delete(store, services.NewBuildingDelete(pid))
				if err != nil {
					Fail(err.Error())
				}
				By("Delete data ok")
			})
		})

		Context("Get 1 record", func() {
			It("should return ok", func() {
				name := fmt.Sprintf("marina-bay-sands::%s", fake.DigitsN(15))
				pid, err := service.Create(store,
					&services.BuildingCreateParams{
						Name:    &name,
						Address: fmt.Sprintf("Marina Boulevard::%s", fake.DigitsN(15)),
						Floors:  tools.Seeder{}.CreateFloors(),
					},
				)
				if err != nil {
					Fail(err.Error())
				}
				Expect(pid).Should(BeNumerically(">", 0))
				By("Create data before get 1 row ok")

				row, err := service.Get(store, services.NewBuildingGetOne(pid))
				if err != nil {
					Fail(err.Error())
				}
				Expect(row.ID).Should(BeNumerically(">", 0))
				By("Get 1 data ok")
			})
		})

		Context("Get list of records", func() {
			It("should return ok", func() {

				for i := 1; i <= 5; i++ {
					name := fmt.Sprintf("marina-bay-sands::%s", fake.DigitsN(15))
					pid, err := service.Create(store,
						&services.BuildingCreateParams{
							Name:    &name,
							Address: fmt.Sprintf("Marina Boulevard::%s", fake.DigitsN(15)),
							Floors:  tools.Seeder{}.CreateFloors(),
						},
					)
					if err != nil {
						Fail(err.Error())
					}
					Expect(pid).Should(BeNumerically(">", 0))
					By("Create data before get list ok")
				}

				rows, err := service.GetAll(
					store,
					&tools.PagingParams{
						Page:  1,
						Limit: 10,
					})
				if err != nil {
					Fail(err.Error())
				}
				Expect(len(rows.Items)).Should(BeNumerically(">", 0))
				Expect(rows.Total).Should(BeNumerically(">", 0))
				By("Get more data ok")
			})
		})

		Context("Create record with minimum parameter", func() {
			It("should return ok", func() {
				name := fmt.Sprintf("building::%s", fake.DigitsN(15))
				pid, err := service.Create(store,
					&services.BuildingCreateParams{
						Name: &name,
					})
				if err != nil {
					Fail(err.Error())
				}
				Expect(pid).Should(BeNumerically(">", 0))
				By("Create data ok")

			})
		})

	}) // valid

	Context("Invalid parameters", func() {

		Context("Update record with missing parameter id", func() {
			It("should error", func() {
				name := fmt.Sprintf("marina-bay-sands::%s", fake.DigitsN(15))
				pid, err := service.Create(store,
					&services.BuildingCreateParams{
						Name:    &name,
						Address: fmt.Sprintf("Marina Boulevard::%s", fake.DigitsN(15)),
						Floors:  tools.Seeder{}.CreateFloors(),
					})
				if err != nil {
					Fail(err.Error())
				}
				Expect(pid).Should(BeNumerically(">", 0))
				By("Create data before update ok")

				fmtPid, _ := strconv.ParseInt(
					fmt.Sprintf("88%d99", int64(rand.Intn(99999999))),
					10, 64)
				uparams := &services.BuildingUpdateParams{
					ID: &fmtPid,
					BuildingCreateParams: services.BuildingCreateParams{
						Name:    &name,
						Address: fmt.Sprintf("Marina Boulevard::%s", fake.DigitsN(15)),
						Floors:  tools.Seeder{}.CreateFloors(),
					},
				}

				err = service.Update(store, uparams)
				Expect(err).To(HaveOccurred())
				By("Update data empty as expected")
			})
		})

		Context("Delete record with missing parameter id", func() {
			It("should error", func() {
				name := fmt.Sprintf("marina-bay-sands::%s", fake.DigitsN(15))
				pid, err := service.Create(store,
					&services.BuildingCreateParams{
						Name:    &name,
						Address: "Marina Boulevard-1a",
						Floors:  tools.Seeder{}.CreateFloors(),
					})
				if err != nil {
					Fail(err.Error())
				}
				Expect(pid).Should(BeNumerically(">", 0))
				By("Create data before delete ok")

				err = service.Delete(store,
					services.NewBuildingDelete(pid+int64(rand.Intn(9999999999))))
				Expect(err).To(HaveOccurred())
				By("Delete data empty as expected")
			})
		})

		Context("Get a record not exists", func() {
			It("should error", func() {
				row, err := service.Get(store,
					&services.BuildingGetParams{
						ID: int64(rand.Intn(9999999999) + rand.Intn(9999999999)),
					})
				Expect(err).To(HaveOccurred())
				Expect(row).To(BeZero())
				By("Get data empty as expected")
			})
		})
	}) // invalid
})
