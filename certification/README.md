# Certification C3 — Oscar Mauricio Areiza Paramo

## Cover

| Field | Value |
|---|---|
| Full name | Oscar Mauricio Areiza Paramo |
| GitHub user | OscarAreiza |
| Group | G2 |
| Team | lms-library |
| Team -docs repository | https://github.com/code-corhuila/library-docs |
| Course | HCCDP-Solution Architectures Course |
| Course link | https://e.huawei.com/en/talent/#/ (Huawei ICT Academy — HCCDP-Solution Architectures Course) |
| Start date | 2026-11-01 |
| Completion date | 2026-11-02 |

Certificate of Completion: `HCCDP-Solution Architectures Course.png` (code ICT20261103000001, issue date 2026-11-03 in the platform's time zone).

## Study log

| Date | Hours | Chapters or labs | Evidence file |
|---|---|---|---|
| 2026-11-01 | 0.5 | Enrollment and pre-class guidance | evidence/00-enrollment.png |
| 2026-11-01 | 1 | Chapter 1 — Architecture Design Methodology on Huawei Cloud | evidence/01-chapter-01.png |
| 2026-11-01 | 1 | Chapter 2 — Typical Service Implementation Solutions on the Cloud | evidence/02-chapter-02.png |
| 2026-11-01 | 1 | Chapter 3 — Cloud Compute Solution Design | evidence/03-chapter-03.png |
| 2026-11-01 | 1 | Chapter 4 — Cloud Storage Solution Design | evidence/04-chapter-04.png |
| 2026-11-01 | 1 | Chapter 5 — Cloud Network Solution Design | evidence/05-chapter-05.png |
| 2026-11-01 | 1 | Chapter 6 — Cloud Database Solution Design | evidence/06-chapter-06.png |
| 2026-11-01 | 1 | Chapter 7 — Cloud Security Solution Design | evidence/07-chapter-07.png |
| 2026-11-01 | 1 | Chapter 8 — Cloud O&M Solution Design | evidence/08-chapter-08.png |
| 2026-11-01 | 1 | Chapter 9 — Distributed Architecture Design | evidence/09-chapter-09.png |
| 2026-11-01 | 1 | Chapter 10 — Highly Scalable System Design | evidence/10-chapter-10.png |
| 2026-11-01 | 1 | Chapter 11 — Disaster Recovery System Design | evidence/11-chapter-11.png |
| 2026-11-01 | 1 | Chapter 12 — Automated Deployment | evidence/12-chapter-12.png |
| 2026-11-01 | 1 | Chapter 13 — Containers and Agility | evidence/13-chapter-13.png |
| 2026-11-01 | 1 | Chapter 14 — Huawei Cloud Advanced Services | evidence/14-chapter-14.png |
| 2026-11-02 | 2.5 | Lab — Design Cases and Practices of an O&M Solution | evidence/15-lab-1.png |
| 2026-11-02 | 2 | Lab — High-Performance Network Solution | evidence/16-lab-2.png |
| 2026-11-02 | 2 | Lab — Highly Elastic Service System | evidence/17-lab-3.png |
| 2026-11-02 | 0.5 | Final exam (mock exam, 18 min 51 s) | evidence/19-final-exam.png |
| Total | 23.5 | | |

## Labs

| Evidence file | Lab name | What I built | What failed first and how I fixed it |
|---|---|---|---|
| evidence/15-lab-1.png | Design Cases and Practices of an O&M Solution | I configured Log Tank Service (LTS) to collect the application logs of an ECS and Cloud Eye alarm rules on CPU and memory that notify through SMN. | The log group showed no data at first because the ICAgent was not installed on the ECS; I installed it from the LTS console and the log stream started receiving entries. |
| evidence/16-lab-2.png | High-Performance Network Solution | I connected two VPCs with a VPC peering connection, gave private subnets Internet egress through a NAT Gateway with an EIP, and resolved the service by name with private DNS. | Ping between the two VPCs failed after creating the peering because the route tables had no route to the peer CIDR; I added the routes on both sides and traffic flowed. |
| evidence/17-lab-3.png | Highly Elastic Service System | I deployed Discuz! behind an ELB with an Auto Scaling group of ECS instances that scales out on CPU usage. | New instances failed the ELB health check because the security group did not allow the health-check port from the ELB; I opened that port and the instances became healthy. |

## Final exam

- Score: 79 points (Passed, 24/30 closed-ended questions correct)
- Date: 2026-11-02
- Evidence: evidence/19-final-exam.png

## Skills summary

I learned to design a cloud solution by checking it against the five principles: reliability, security, performance efficiency, cost and O&M.
I can make stateless services scale horizontally with ELB plus Auto Scaling policies driven by a concrete metric and threshold.
I can design cross-AZ high availability and choose a DR pattern (backup and restore, active-standby, active-active) from the RTO and RPO a data store needs.
I can design the network of a solution with VPCs, subnets, security groups, VPC peering, NAT Gateway and private DNS.
I can plan observability with LTS for logs and Cloud Eye for metrics and alarms, and automate infrastructure with RFS templates.
I also understand how containers on CCE and a CI/CD pipeline make microservice deployment repeatable and easy to roll back.
