var config = {
    endpoint: "http://dbpedia.org/sparql"
}

var formatQueryForAutocompletion = function(yasqe, partialToken, query) {
     var cur = yasqe.getCursor(false);
     var begin = yasqe.getRange({line: 0, ch:0}, cur);
     query = begin + "< " + query.substring(begin.length, query.length);
     return query;
};

/**
 * Sparqled autocompletion function
 */

var customAutocompletionFunction = function(yasqe, partialToken, type, callback) {
    autocompletion.RecommendationQuery(formatQueryForAutocompletion(yasqe, partialToken, yasqe.getValue()), function(q, err) {
        if (err) {
            alert(err)
            return
        }
        if (!q) {
            alert("No recommendation at this position")
            return
        }
        var ajaxConfig = {
            type: "GET",
            crossDomain: true,
            url: config.endpoint,
            data: {
                format: 'application/json',
                query: q
            },
            success: function(data) {
                var completions = [];
                for (var i = 0; i < data.results.bindings.length; i++) {
                    var binding = data.results.bindings[i];
                    var pof = binding.POF.value
                    // The YASQE library automatically wraps the string with '<' and '>'
                    completions.push(pof.substring(1, pof.length - 1));
                }
                callback(completions);
            },
            beforeSend: function(){
                $('#loading').show();
            },
            complete: function(){
                $('#loading').hide();
            }
        };
        $.ajax(ajaxConfig);
    })
};

var yasqe = YASQE(document.getElementById("yasqe"), {
	sparql: {
        endpoint: config.endpoint,
		showQueryButton: true,
	},
	autocompletions: {
		classes: {
			async: true,
			get: customAutocompletionFunction
		},
		properties: {
			async: true,
			get: customAutocompletionFunction
		}
	}
});
var yasr = YASR(document.getElementById("yasr"), {
	getUsedPrefixes: yasqe.getPrefixesFromQuery
});

/**
* Set some of the hooks to link YASR and YASQE
*/
yasqe.options.sparql.handlers.success =  function(data, status, response) {
	yasr.setResponse({response: data, contentType: response.getResponseHeader("Content-Type")});
};
yasqe.options.sparql.handlers.error = function(xhr, textStatus, errorThrown) {
	yasr.setResponse({exception: textStatus + ": " + errorThrown});
};
