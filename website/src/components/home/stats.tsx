import { component$ } from "@builder.io/qwik";

const stats = [
  { value: "99.9%", label: "Uptime" },
  { value: "<100ms", label: "API Response" },
  { value: "3", label: "Official SDKs" },
  { value: "5min", label: "Setup Time" },
];

export const Stats = component$(() => {
  return (
    <div class="stats-row">
      {stats.map((stat) => (
        <div class="stats-row__item" key={stat.label}>
          <span class="stats-row__value">{stat.value}</span>
          <span class="stats-row__label">{stat.label}</span>
        </div>
      ))}
    </div>
  );
});
