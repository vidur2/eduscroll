context = """
Series 267

Despite the previous example, in practice definitions 13.7 and 13.5 are
rarely used to prove that a particular sequence or series converges to a
particular number. Instead we tend to use a multitude of convergence tests
that are covered in a typical calculus course. Examples of such tests include
the comparison test, the ratio test, the root test and the alternating series test.
You learned how to use these tests and techniques in your calculus course,
though that course may not have actually proved that the tests were valid.
The point of our present discussion is that definitions 13.7 and 13.5 can be
used to prove the tests. To underscore this point, this section’s exercises
ask you to prove several convergence tests.
By way of illustration, we close with a proof of a theorem that leads to a
test for divergence.

Theorem 13.11 If

∞
∑

k=1

{ }

a k converges, then the sequence a n converges to 0.

Proof. We use direct proof. Suppose

∞
∑

k=1

a k converges, and say

{ }

∞
∑

k=1

a k = S.

Then by Definition 13.7, the sequence of partial sums s n converges to S.
From this, Definition
∣ ∣ 13.5 says that for any ε > 0 there∣ is an N ∣∈ N for which
n > N implies ∣ s n − S ∣ < ε. Thus also n − 1 > N implies ∣ s n−1 − S ∣ < ε.
{ }
We need to show that a n converges to 0. So take ε > 0. ∣By the∣ previous
paragraph, there is
an N ′ ∈ N for which n > N ′ implies ∣ s n − S ∣ < 2ε and
∣
∣ ε
∣ s n−1 − S ∣ < . Notice that a n = s n − s n−1 for any n > 2. So if n > N ′ we have
2

∣ ∣ ∣
∣ ∣( ) ( )∣
∣ a n − 0∣ = ∣ s n − s n−1 ∣ = ∣ s n − S − s n−1 − S ∣
∣ ∣ ∣ ∣
ε ε
≤ ∣ s n − S ∣ + ∣ s n−1 − S ∣ < + = ε.
2 2
{ }
Therefore, by Definition 13.5, the sequence a n converges to 0.

■

The contrapositive of this theorem is a convenient test for divergence:

{ }

Corollary 13.1 (Divergence test) If a n diverges, or if it converges to a
∞
∑
non-zero number, then
a k diverges.

k=1

For example, according to the divergence test, the series

{ }

∞
∑

∞ ( )
∑
1 − 1k

k=1

k+1

(−1) (k+1)
diverges, because the sequence 1 − n1 converges to 1. Also, k
k=1
{ }
n+1
diverges because (−1) n (n+1) diverges. (See Example 13.9 on page 263.)
The divergence test gives only a criterion for deciding if a series diverges.

{ }

It says nothing about convergence. If a n converges to 0, then

∞
∑

k=1

Free PDF version

a k may
"""