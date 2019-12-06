The current policy filter is quite expensive (10 seconds). Seems I am missing an index somewhere. 



```sql
SELECT * FROM objects;
                                                                                              QUERY PLAN                                                                                              
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
 Limit  (cost=0.00..4833092.71 rows=1 width=24) (actual time=9421.505..9421.509 rows=1 loops=1)
   ->  Seq Scan on objects  (cost=0.00..387841190690.00 rows=80247 width=24) (actual time=9421.503..9421.503 rows=1 loops=1)
         Filter: (((SubPlan 1) >= 'read'::permission_level) OR ((SubPlan 3) >= 'read'::permission_level) OR ((SubPlan 4) >= 'owner'::permission_level) OR ((SubPlan 6) >= 'owner'::permission_level))
         Rows Removed by Filter: 29
         SubPlan 1
           ->  Index Scan using object_user_permissions_pkey on object_user_permissions  (cost=0.43..8.45 rows=1 width=4) (actual time=0.009..0.009 rows=0 loops=30)
                 Index Cond: ((user_id = (current_setting('jwt.userid'::text))::integer) AND (object_id = objects.id))
         SubPlan 3
           ->  Seq Scan on object_group_permissions  (cost=0.00..1939197.49 rows=224423 width=4) (actual time=154.422..159.961 rows=0 loops=30)
                 Filter: (SubPlan 2)
                 Rows Removed by Filter: 448846
                 SubPlan 2
                   ->  Result  (cost=0.29..8.31 rows=1 width=4) (actual time=0.000..0.000 rows=0 loops=13465380)
                         One-Time Filter: (object_group_permissions.object_id = objects.id)
                         ->  Index Only Scan using users_groups_pkey on users_groups  (cost=0.29..8.31 rows=1 width=4) (actual time=0.002..0.002 rows=1 loops=131)
                               Index Cond: (user_id = (current_setting('jwt.userid'::text))::integer)
                               Heap Fetches: 131
         SubPlan 4
           ->  Index Scan using object_user_permissions_pkey on object_user_permissions object_user_permissions_1  (cost=0.43..8.45 rows=1 width=4) (actual time=0.009..0.009 rows=0 loops=29)
                 Index Cond: ((user_id = (current_setting('jwt.userid'::text))::integer) AND (object_id = objects.id))
         SubPlan 6
           ->  Seq Scan on object_group_permissions object_group_permissions_1  (cost=0.00..1939197.49 rows=224423 width=4) (actual time=159.366..159.366 rows=0 loops=29)
                 Filter: (SubPlan 5)
                 Rows Removed by Filter: 448846
                 SubPlan 5
                   ->  Result  (cost=0.29..8.31 rows=1 width=4) (actual time=0.000..0.000 rows=0 loops=13016534)
                         One-Time Filter: (object_group_permissions_1.object_id = objects.id)
                         ->  Index Only Scan using users_groups_pkey on users_groups users_groups_1  (cost=0.29..8.31 rows=1 width=4) (actual time=0.002..0.002 rows=1 loops=124)
                               Index Cond: (user_id = (current_setting('jwt.userid'::text))::integer)
                               Heap Fetches: 124
 Planning time: 1.265 ms
 Execution time: 9421.668 ms
(32 rows)
```