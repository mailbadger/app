<?php

namespace newsletters\Entities;

use Illuminate\Database\Eloquent\Model;
use Prettus\Repository\Contracts\Transformable;
use Prettus\Repository\Traits\TransformableTrait;

class Field extends Model implements Transformable
{
    use TransformableTrait;

    protected $table = 'fields';

    protected $fillable = [
        'name',
        'list_id'
    ];

    public function subList()
    {
        return $this->belongsTo('newsletters\Entities\List', 'list_id');
    }
}
