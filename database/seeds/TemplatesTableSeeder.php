<?php

use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;

class TemplatesTableSeeder extends Seeder
{
    /**
     * Run the database seeds.
     *
     * @return void
     */
    public function run()
    {
        for ($i =0; $i < 10; $i++) {
            DB::table('templates')->insert([
                'name' => 'Template '.$i,
                'content' => 'Dummy Content'
            ]);
        }
    }
}
