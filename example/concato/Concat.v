Inductive linkedList : Set :=
  | Nil : linkedList
  | Prepend : nat -> linkedList -> linkedList.

Inductive concato (xs ys zs: linkedList): Prop :=
  | concatNil:
     xs = Nil ->
     zs = ys ->
     concato xs ys zs
  | concatPrepend:
     (exists head xtail ztail,
        xs = Prepend head xtail /\
        zs = Prepend head ztail /\
        concato xtail ys ztail) ->
        concato xs ys zs.

Theorem concatNil_is_correct:
  forall xs ys zs,
    xs = Nil ->
    zs = ys ->
    concato xs ys zs.
Proof.
  intros.
  rewrite H.
  rewrite H0.
  apply concatNil.
  reflexivity.
  reflexivity.
Qed.

Theorem concatPrepend_is_correct:
  forall xs ys zs head xtail ztail,
    xs = Prepend head xtail ->
    zs = Prepend head ztail ->
    concato xtail ys ztail ->
    concato xs ys zs.
Proof.
  intros.
  rewrite H.
  rewrite H0.
  apply concatPrepend.
  exists head.
  exists xtail.
  exists ztail.
  split.
  reflexivity.
  split.
  reflexivity.
  apply H1.
Qed.

Fixpoint concat (xs ys: linkedList): linkedList :=
  match xs with
  | Nil => ys
  | Prepend head xtail => Prepend head (concat xtail ys)
  end.

Theorem concat_is_correct:
  forall xs ys zs,
    concato xs ys zs <-> concat xs ys = zs.
Proof.
  split; generalize dependent zs; generalize dependent ys.
  - induction xs; intros.
    + simpl.
      destruct H.
      * subst. reflexivity.
      * destruct H as [head [xtail [ztail [H1 [H2 H3]]]]].
        discriminate.
    + simpl.
      destruct H.
      * discriminate.
      * destruct H as [head [xtail [ztail [H1 [H2 H3]]]]].
        subst.
        inversion H1.
        subst.
        rewrite (IHxs ys ztail); try reflexivity.
        assumption.
  - induction xs.
    + intros.
      rewrite <- H.
      apply concatNil.
      reflexivity.
      reflexivity.
    + intros.
      rewrite <- H.
      apply concatPrepend.
      exists n.
      exists xs.
      exists (concat xs ys).
      split.
      reflexivity.
      split.
      reflexivity.
      apply IHxs.
      reflexivity.
Qed.



