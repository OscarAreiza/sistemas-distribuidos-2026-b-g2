# ADR-C3-01: Run the stateless LMS services on CCE behind ELB with horizontal autoscaling

## Status

Proposed

## Context

Today every LMS service (`lms-catalog-api`, `lms-membership-api`, `lms-circulation-api`, `lms-access-api`, `lms-workflow`, the portals and `lms-api-gateway`) runs as a container on one Docker host assembled with `lms-pod/compose.yml` (`12-week/hu-status/deployment-plan.md`, section 2). That host is a single point of failure, and the only way to absorb a load peak (for example, the start of the semester, when most loans are registered) is to resize the whole machine by hand. The design note names this as the reliability and performance weakness in section 2 (Five principles) and sets per-service scaling metrics in section 3 (Scalability). The services are already stateless and packaged as Docker images (`05-week/hu-status/CONTAINERIZATION.md`), and every state lives in the data stores, so they can be replicated without code changes.

## Decision

We will deploy the stateless services as Deployments on a Cloud Container Engine (CCE) cluster whose worker nodes are spread across two availability zones, pull their images from SWR (SoftWare Repository for Container), and expose them through a shared ELB that replaces the single-host entry point in front of `lms-api-gateway`. Each Deployment gets a Horizontal Pod Autoscaler with the metric, threshold and minimum/maximum replicas listed in section 3 of the design note (minimum 2 replicas, so one pod per AZ survives an AZ loss), and the CCE cluster autoscaler adds nodes when pods cannot be scheduled. Readiness and liveness probes reuse the `/health` endpoints the rolling deployment already checks.

## Alternatives considered

- Keep a single ECS with Docker Compose: it is the cheapest option and what we run today, but it keeps the single point of failure and cannot scale one service independently of the others.
- ECS instances in an Auto Scaling group behind ELB (one VM image per service): it gives elasticity and cross-AZ placement, but it scales whole VMs per service, which wastes capacity for 12 small services and makes image rollout and rollback slower than a container rollout.

## Consequences

- Positive: no single host failure takes the system down, each service scales on its own metric, and rolling updates and rollbacks become a `kubectl rollout` on an immutable SWR image tag instead of a manual script.
- Negative: the team must learn and operate Kubernetes objects (Deployments, Services, HPA, probes), and the minimum of two replicas per service plus the cluster nodes costs more than the single host, which goes against our minimal-infrastructure principle (P5) and must be justified with the reliability gain.
