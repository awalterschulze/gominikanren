same(X, X) ->
    true;
same(X, Y) ->
    false.

compress([]) ->
    [];
compress([X, X | Rest]) ->
    compress([X | Rest]);
compress([X | Rest]) ->
    [X | compress(Rest)].