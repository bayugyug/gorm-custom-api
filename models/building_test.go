package models_test

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/bayugyug/gorm-custom-api/api/routes"
	"github.com/bayugyug/gorm-custom-api/models"
	"github.com/bayugyug/gorm-custom-api/tools"
	"github.com/icrowley/fake"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var svcrouter *routes.APIRouter
var store *gorm.DB

var _ = BeforeSuite(func() {
	var err error
	//init
	svcrouter, err = routes.NewAPIRouter(
		routes.WithSvcOptConfig(tools.DevelTester{}.Config()),
	)
	Expect(err).NotTo(HaveOccurred())

})

var _ = AfterSuite(func() {
	tools.DevelTester{}.EmptyDBTst(svcrouter.Building.Storage)
})

var _ = Describe("REST Building API Service::MODELS", func() {

	BeforeEach(func() {
		store = svcrouter.Building.Storage
	})

	Context("Valid parameters", func() {

		Context("Create record", func() {
			It("should return ok", func() {
				name := fmt.Sprintf("building::%s", fake.DigitsN(15))
				params := &models.BuildingCreateParams{
					Name:    &name,
					Address: fmt.Sprintf("Marina Boulevard::%s", fake.DigitsN(15)),
					Floors:  tools.Seeder{}.CreateFloors(),
				}
				pid, err := params.Create(store)
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
				params := &models.BuildingCreateParams{
					Name:    &name,
					Address: fmt.Sprintf("Marina Boulevard::%s", fake.DigitsN(15)),
					Floors:  tools.Seeder{}.CreateFloors(),
				}
				pid, err := params.Create(store)
				if err != nil {
					Fail(err.Error())
				}
				Expect(pid).Should(BeNumerically(">", 0))
				By("Create data before update ok")

				uparams := &models.BuildingUpdateParams{
					ID: &pid,
					BuildingCreateParams: models.BuildingCreateParams{
						Name:    &name,
						Address: fmt.Sprintf("Marina Boulevard::%s", fake.DigitsN(15)),
						Floors:  tools.Seeder{}.CreateFloors(),
					},
				}

				err = uparams.Update(store)
				if err != nil {
					Fail(err.Error())
				}
				By("Update data ok")
			})
		})

		Context("Delete record", func() {
			It("should return ok", func() {
				name := fmt.Sprintf("marina-bay-sands::%s", fake.DigitsN(15))
				params := &models.BuildingCreateParams{
					Name:    &name,
					Address: fmt.Sprintf("Marina Boulevard::%s", fake.DigitsN(15)),
					Floors:  tools.Seeder{}.CreateFloors(),
				}
				pid, err := params.Create(store)
				if err != nil {
					Fail(err.Error())
				}
				Expect(pid).Should(BeNumerically(">", 0))
				By("Create data before delete ok")

				uparams := models.NewBuildingDelete(pid)
				err = uparams.Delete(store)
				if err != nil {
					Fail(err.Error())
				}
				By("Delete data ok")
			})
		})

		Context("Get 1 record", func() {
			It("should return ok", func() {
				name := fmt.Sprintf("marina-bay-sands::%s", fake.DigitsN(15))
				params := &models.BuildingCreateParams{
					Name:    &name,
					Address: fmt.Sprintf("Marina Boulevard::%s", fake.DigitsN(15)),
					Floors:  tools.Seeder{}.CreateFloors(),
				}
				pid, err := params.Create(store)
				if err != nil {
					Fail(err.Error())
				}
				Expect(pid).Should(BeNumerically(">", 0))
				By("Create data before get 1 row ok")

				uparams := models.NewBuildingGetOne(pid)
				row, err := uparams.Get(store)
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
					params := &models.BuildingCreateParams{
						Name:    &name,
						Address: fmt.Sprintf("Marina Boulevard::%s", fake.DigitsN(15)),
						Floors:  tools.Seeder{}.CreateFloors(),
					}
					pid, err := params.Create(store)
					if err != nil {
						Fail(err.Error())
					}
					Expect(pid).Should(BeNumerically(">", 0))
					By("Create data before get list ok")
				}
				uparams := &models.BuildingGetParams{}
				rows, err := uparams.GetAll(store)
				if err != nil {
					Fail(err.Error())
				}
				Expect(len(rows)).Should(BeNumerically(">", 0))
				By("Get more data ok")
			})
		})

		Context("Create record with minimum parameter", func() {
			It("should return ok", func() {
				name := fmt.Sprintf("building::%s", fake.DigitsN(15))
				params := &models.BuildingCreateParams{
					Name: &name,
				}
				pid, err := params.Create(store)
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
				params := &models.BuildingCreateParams{
					Name:    &name,
					Address: fmt.Sprintf("Marina Boulevard::%s", fake.DigitsN(15)),
					Floors:  tools.Seeder{}.CreateFloors(),
				}

				pid, err := params.Create(store)
				if err != nil {
					Fail(err.Error())
				}
				Expect(pid).Should(BeNumerically(">", 0))
				By("Create data before update ok")
				fmtPid, _ := strconv.ParseInt(
					fmt.Sprintf("/v1/api/building/999%d888",
						int64(rand.Intn(9999999999))),
					10, 64)
				uparams := &models.BuildingUpdateParams{
					ID: &fmtPid,
					BuildingCreateParams: models.BuildingCreateParams{
						Name:    &name,
						Address: fmt.Sprintf("Marina Boulevard::%s", fake.DigitsN(15)),
						Floors:  tools.Seeder{}.CreateFloors(),
					},
				}

				err = uparams.Update(store)
				Expect(err).To(HaveOccurred())
				By("Update data empty as expected")
			})
		})

		Context("Delete record with missing parameter id", func() {
			It("should error", func() {
				name := fmt.Sprintf("marina-bay-sands::%s", fake.DigitsN(15))
				params := &models.BuildingCreateParams{
					Name:    &name,
					Address: "Marina Boulevard-1a",
					Floors:  tools.Seeder{}.CreateFloors(),
				}
				pid, err := params.Create(store)
				if err != nil {
					Fail(err.Error())
				}
				Expect(pid).Should(BeNumerically(">", 0))
				By("Create data before delete ok")

				uparams := models.NewBuildingDelete(pid + int64(rand.Intn(9999999999)))
				err = uparams.Delete(store)
				Expect(err).To(HaveOccurred())
				By("Delete data empty as expected")
			})
		})

		Context("Get a record not exists", func() {
			It("should error", func() {
				uparams := &models.BuildingGetParams{ID: int64(rand.Intn(9999999999) + rand.Intn(9999999999))}
				row, err := uparams.Get(store)
				Expect(err).To(HaveOccurred())
				Expect(row).To(BeZero())
				By("Get data empty as expected")
			})
		})
	}) // invalid
})
