package shadow // import "github.com/kihamo/shadow"

//go:generate goimports -w ./
//go:generate sh -c "cd components/alerts && go-bindata-assetfs -pkg=alerts templates/..."
//go:generate sh -c "cd components/dashboard && go-bindata-assetfs -pkg=dashboard templates/... public/..."
//go:generate sh -c "cd components/mail && go-bindata-assetfs -pkg=mail templates/..."
//go:generate sh -c "cd components/workers && go-bindata-assetfs -pkg=workers templates/..."
