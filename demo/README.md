Gosparqled library as a plugin for YASQE
========================================

[YASQE](http://yasqe.yasgui.org/) provides a plugin mechanism for providing autocompletions.

The steps taken for returning recommendations are the following:
- Gosparqled processes the input query and creates a new one that projects the **POF**. The POF is the position in the input query that is to be completed;
- the created query is executed against the endpoint using an AJAX call; and
- the bindings of the POF are returned as recommendations to the user.
