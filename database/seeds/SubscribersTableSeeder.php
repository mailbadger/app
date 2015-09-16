<?php

use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;

class SubscribersTableSeeder extends Seeder
{
    /**
     * Run the database seeds.
     *
     * @return void
     */
    public function run()
    {
        DB::table('subscribers')->delete();
        factory(newsletters\Entities\Subscriber::class, 100)->create()->each(function ($s) {
            //$s->lists()->attach(5);
        });
    }
}
