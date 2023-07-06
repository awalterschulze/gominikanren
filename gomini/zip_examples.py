def zip(xs: list[A], ys: list[B]): list[tuple[A, B]] =
    return ((x, y) for x in xs for y in ys)

xs = [1, 2, 3]
ys = ["a", "b", "c"]
zs = zip(xs, ys)
print(tuple(zs))
zs = [(1, "a"), (2, "b"), (3, "c")]
# Output: [(1, 'a'), (2, 'b'), (3, 'c')]

xs = [1, 2]
ys = ["a", "b", "c”, “d”]
zs = zip(xs, ys)
zs = None


def reduce(xs: list[A], init: B, f: Callable[[A, B], B]): B =
    acc = init

xs = [1, 2, 3]
acc = reduce(xs, 0, lambda x, accumulator: x + accumulator)
acc = 6

xs = [True, False, True]
acc = reduce(xs, True, lambda x, accumulator: x and accumulator)
acc = False

xs = [1, 2, 3]
acc = reduce(xs, 1, lambda x, acc: str(x) + acc)
acc = "123"

xs = [1,2,3]
acc = reduce(xs, "", lambda x, acc: acc + str(x))
xs = "123"