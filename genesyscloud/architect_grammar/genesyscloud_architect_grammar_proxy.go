package architect_grammar

import (
	"context"
	"github.com/mypurecloud/platform-client-sdk-go/v109/platformclientv2"
)

/*
The genesyscloud_architect_grammar_proxy.go file contains the proxy structures and methods that interact
with the Genesys Cloud SDK. We use composition here for each function on the proxy so individual functions can be stubbed
out during testing.
*/

// internalProxy holds a proxy instance that can be used throughout the package
var internalProxy *architectGrammarProxy

// Type definitions for each func on our proxy so we can easily mock them out later
type createArchitectGrammarFunc func(ctx context.Context, p *architectGrammarProxy, grammar *platformclientv2.Grammar) (*platformclientv2.Grammar, error)
type getAllArchitectGrammarFunc func(ctx context.Context, p *architectGrammarProxy) (*[]platformclientv2.Grammar, error)
type getArchitectGrammarByIdFunc func(ctx context.Context, p *architectGrammarProxy, grammarId string) (grammar *platformclientv2.Grammar, responseCode int, err error)
type getArchitectGrammarIdByNameFunc func(ctx context.Context, p *architectGrammarProxy, name string) (grammarId string, retryable bool, err error)
type updateArchitectGrammarFunc func(ctx context.Context, p *architectGrammarProxy, grammarId string, grammar *platformclientv2.Grammar) (*platformclientv2.Grammar, error)
type deleteArchitectGrammarFunc func(ctx context.Context, p *architectGrammarProxy, grammarId string) (responseCode int, err error)

// architectGrammarProxy contains all of the methods that call genesys cloud APIs.
type architectGrammarProxy struct {
	clientConfig                    *platformclientv2.Configuration
	architectApi                    *platformclientv2.ArchitectApi
	createArchitectGrammarAttr      createArchitectGrammarFunc
	getAllArchitectGrammarAttr      getAllArchitectGrammarFunc
	getArchitectGrammarByIdAttr     getArchitectGrammarByIdFunc
	getArchitectGrammarIdByNameAttr getArchitectGrammarIdByNameFunc
	updateArchitectGrammarAttr      updateArchitectGrammarFunc
	deleteArchitectGrammarAttr      deleteArchitectGrammarFunc
}

// newArchitectGrammarProxy initializes the grammar proxy with all of the data needed to communicate with Genesys Cloud
func newArchitectGrammarProxy(clientConfig *platformclientv2.Configuration) *architectGrammarProxy {
	api := platformclientv2.NewArchitectApiWithConfig(clientConfig)
	return &architectGrammarProxy{
		clientConfig:                    clientConfig,
		architectApi:                    api,
		createArchitectGrammarAttr:      createArchitectGrammarFn,
		getAllArchitectGrammarAttr:      getAllArchitectGrammarFn,
		getArchitectGrammarByIdAttr:     getArchitectGrammarByIdFn,
		getArchitectGrammarIdByNameAttr: getArchitectGrammarIdByNameFn,
		updateArchitectGrammarAttr:      updateArchitectGrammarFn,
		deleteArchitectGrammarAttr:      deleteArchitectGrammarFn,
	}
}

// getArchitectGrammarProxy acts as a singleton for the internalProxy. It also ensures
// that we can still proxy our tests by directly setting internalProxy package variable
func getArchitectGrammarProxy(clientConfig *platformclientv2.Configuration) *architectGrammarProxy {
	if internalProxy == nil {
		internalProxy = newArchitectGrammarProxy(clientConfig)
	}

	return internalProxy
}

// createArchitectGrammar creates a Genesys Cloud Architect Grammar
func (p *architectGrammarProxy) createArchitectGrammar(ctx context.Context, grammar *platformclientv2.Grammar) (*platformclientv2.Grammar, error) {
	return p.createArchitectGrammarAttr(ctx, p, grammar)
}

// getAllArchitectGrammar retrieves all Genesys Cloud Architect Grammar
func (p *architectGrammarProxy) getAllArchitectGrammar(ctx context.Context) (*[]platformclientv2.Grammar, error) {
	return p.getAllArchitectGrammarAttr(ctx, p)
}

// getArchitectGrammarById returns a single Genesys Cloud Architect Grammar by Id
func (p *architectGrammarProxy) getArchitectGrammarById(ctx context.Context, grammarId string) (grammar *platformclientv2.Grammar, statusCode int, err error) {
	return p.getArchitectGrammarByIdAttr(ctx, p, grammarId)
}

// getArchitectGrammarIdByName returns a single Genesys Cloud Architect Grammar by a name
func (p *architectGrammarProxy) getArchitectGrammarIdByName(ctx context.Context, name string) (grammarId string, retryable bool, err error) {
	return p.getArchitectGrammarIdByNameAttr(ctx, p, name)
}

// updateArchitectGrammar updates a Genesys Cloud Architect Grammar
func (p *architectGrammarProxy) updateArchitectGrammar(ctx context.Context, grammarId string, grammar *platformclientv2.Grammar) (*platformclientv2.Grammar, error) {
	return p.updateArchitectGrammarAttr(ctx, p, grammarId, grammar)
}

// deleteArchitectGrammar deletes a Genesys Cloud Architect Grammar by Id
func (p *architectGrammarProxy) deleteArchitectGrammar(ctx context.Context, grammarId string) (statusCode int, err error) {
	return p.deleteArchitectGrammarAttr(ctx, p, grammarId)
}

// createArchitectGrammarFn is an implementation function for creating a Genesys Cloud Architect Grammar
func createArchitectGrammarFn(ctx context.Context, p *architectGrammarProxy, grammar *platformclientv2.Grammar) (*platformclientv2.Grammar, error) {
	return nil, nil
}

// getAllArchitectGrammarFn is the implementation for retrieving all Architect Grammars in Genesys Cloud
func getAllArchitectGrammarFn(ctx context.Context, p *architectGrammarProxy) (*[]platformclientv2.Grammar, error) {
	return nil, nil
}

// getArchitectGrammarByIdFn is an implementation of the function to get a Genesys Cloud Architect Grammar by Id
func getArchitectGrammarByIdFn(ctx context.Context, p *architectGrammarProxy, grammarId string) (grammar *platformclientv2.Grammar, statusCode int, err error) {
	return nil, 0, nil
}

// getArchitectGrammarIdBySearchFn is an implementation of the function to get a Genesys Cloud Architect Grammar by name
func getArchitectGrammarIdByNameFn(ctx context.Context, p *architectGrammarProxy, name string) (grammarId string, retryable bool, err error) {
	return "", false, nil
}

// updateArchitectGrammarFn is an implementation of the function to update a Genesys Cloud Architect Grammar
func updateArchitectGrammarFn(ctx context.Context, p *architectGrammarProxy, grammarId string, grammar *platformclientv2.Grammar) (*platformclientv2.Grammar, error) {
	return nil, nil
}

// deleteArchitectGrammarFn is an implementation function for deleting a Genesys Cloud Architect Grammar
func deleteArchitectGrammarFn(ctx context.Context, p *architectGrammarProxy, grammarId string) (statusCode int, err error) {
	return 0, nil
}
