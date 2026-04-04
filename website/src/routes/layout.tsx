import { component$, Slot } from "@builder.io/qwik";
import { Header } from "../components/layout/header";
import { Footer } from "../components/layout/footer";

export default component$(() => {
  return (
    <>
      <Header />
      <main>
        <Slot />
      </main>
      <Footer />
    </>
  );
});
