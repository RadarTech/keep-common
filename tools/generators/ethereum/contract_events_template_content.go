package main

// contractEventsTemplateContent contains the template string from contract_events.go.tmpl
var contractEventsTemplateContent = `{{- $contract := . -}}
{{- range $i, $event := .Events }}

type {{$contract.FullVar}}{{$event.CapsName}}Func func(
    {{$event.ParamDeclarations -}}
)

func ({{$contract.ShortVar}} *{{$contract.Class}}) Watch{{$event.CapsName}}(
	success {{$contract.FullVar}}{{$event.CapsName}}Func,
	fail func(err error) error,
	{{$event.IndexedFilterDeclarations -}}
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.{{$contract.AbiClass}}{{$event.CapsName}})
	eventSubscription, err := {{$contract.ShortVar}}.contract.Watch{{$event.CapsName}}(
		nil,
		eventChan,
		{{$event.IndexedFilters}}
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for {{$event.CapsName}} events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
                    {{$event.ParamExtractors}}
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

{{- end -}}`
